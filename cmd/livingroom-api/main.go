package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	ggrpc "google.golang.org/grpc"

	"github.com/go-logr/logr"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/version"
	"github.com/sebastianrosch/livingroompresentations/internal/config"
	"github.com/sebastianrosch/livingroompresentations/internal/grpc"
	"github.com/sebastianrosch/livingroompresentations/internal/rest"
	"github.com/sebastianrosch/livingroompresentations/internal/service"
	servicev1 "github.com/sebastianrosch/livingroompresentations/internal/service/v1"
	"github.com/sebastianrosch/livingroompresentations/internal/store"
	"github.com/sebastianrosch/livingroompresentations/pkg/auth"
	"github.com/sebastianrosch/livingroompresentations/pkg/log"
)

func main() {
	logger := log.NewDefaultLogger()
	logger.Info("Starting Living Room API",
		"version", version.Info(),
		"build_context", version.BuildContext())

	// Get the config.
	var httpPort, grpcPort, host string = config.Get().HTTPPort, config.Get().GRPCPort, config.Get().Host

	s, err := store.NewMongoStore(
		config.Get().DatabaseURI,
		config.Get().DatabaseName,
		config.Get().DatabaseUsername,
		config.Get().DatabasePassword,
	)
	if err != nil {
		logger.Error(err, "couldn't create MongoDB store")
		os.Exit(2)
	}

	// Setup the context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.CreateEvent("How viruses spread", "Epidemologist talks about how viruses spread")

	events, _ := s.GetAllEvents()
	for _, event := range events {
		fmt.Print(event)
	}

	httpClient := getHTTPClient()

	tokenDecoder := auth.NewJWTTokenDecoder(config.Get().AuthJwksURL)
	userInfoRetriever := auth.NewUserInfoRetriever(config.Get().AuthUserInfoEndpoint, httpClient)
	metadata := service.NewMetadata()
	whitelist := []string{"/v1.LivingRonom/GetVersion"}
	authContext := &auth.BearerTokenContext{}
	authenticator := auth.NewAuthenticator(logger, whitelist, tokenDecoder, userInfoRetriever, metadata, authContext)

	// Configure the service implementation.
	v1Service := &servicev1.LivingRoomService{}

	// Set up a router to host all handlers on the same port.
	router := setupRouter(ctx, logger, host, grpcPort)

	// Start the HTTP server.
	httpServer := startHTTPServer(logger, host, httpPort, router)

	// Start the gRPC server.
	grpcServer := startgRPCServer(ctx, logger, host, grpcPort, v1Service, authenticator)
	if grpcServer == nil {
		return
	}

	// Setting up signal capturing.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	// Graceful shutdown.
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error(err, "failed shutting down server")
	}
	grpcServer.GracefulStop()

	// Done.
	logger.Info("shutting down")
}

func setupRouter(
	ctx context.Context,
	logger logr.Logger,
	host string,
	grpcPort string) *mux.Router {
	swaggerv1 := packr.New("swagger", "../../api/swagger/v1")
	swaggerui := packr.New("swaggerui", "../../swaggerui")
	schemasv1 := packr.New("schemas", "../../api/schema/v1")

	docsRouter := mux.NewRouter()
	docsRouter.PathPrefix("/docs/swagger/").Handler(http.StripPrefix("/docs/swagger/", http.FileServer(swaggerv1)))
	docsRouter.PathPrefix("/docs/schema/v1/").Handler(http.StripPrefix("/docs/schema/v1/", http.FileServer(schemasv1)))
	docsRouter.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(swaggerui)))

	router := mux.NewRouter()
	router.PathPrefix("/docs/").Handler(docsRouter)
	router.PathPrefix("/api/").Handler(http.StripPrefix("/api", rest.GetHandler(ctx, logger, host, grpcPort)))
	// router.PathPrefix("/").Handler(twirp.GetHandler(ctx, logger, v1Service, authenticator))

	return router
}

// func startHealthCheckEndPoint(logger logr.Logger) {
// 	healthCheckAddr := net.JoinHostPort(config.Get().Host, config.Get().HealthCheckPort)
// 	logger.Info("starting healthcheck", "addr", healthCheckAddr)
// 	go http.ListenAndServe(healthCheckAddr, healthcheck.Handler())
// }

func startHTTPServer(logger logr.Logger, host string, httpPort string, router *mux.Router) *http.Server {
	httpServer := &http.Server{Addr: host + ":" + httpPort, Handler: router}
	go func() {
		logger.Info("starting HTTP server", "addr", host+":"+httpPort)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error(err, "failed to start HTTP server")
			return
		}
	}()

	return httpServer
}

func startgRPCServer(
	ctx context.Context,
	logger logr.Logger,
	host, grpcPort string,
	v1Service *servicev1.LivingRoomService,
	authenticator *auth.TokenAuthenticator) *ggrpc.Server {
	// Only the gRPC server needs a different port, but as it is only internal it doesn't matter.
	listener, err := net.Listen("tcp", host+":"+grpcPort)
	if err != nil {
		logger.Error(err, "failed to bind gRPC server")
		return nil
	}
	grpcServer := grpc.GetServer(ctx, logger, v1Service, authenticator)
	go func() {
		logger.Info("starting gRPC server", "addr", host+":"+grpcPort)
		err := grpcServer.Serve(listener)
		if err != nil {
			logger.Error(err, "failed to start gRPC server")
			return
		}
	}()

	return grpcServer
}

// getHTTPClient returns the HTTP Client instance used in the API.
func getHTTPClient() *http.Client {
	// extracted from https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}

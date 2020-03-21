package rest

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	v1 "github.com/sebastianrosch/livingroompresentations/rpc/livingroom-api/v1"
)

// GetHandler returns the HTTP/REST gateway handler.
func GetHandler(ctx context.Context, logger logr.Logger, host, grpcPort string) http.Handler {
	// Register JSON and YAML marshaler.
	json := &runtime.JSONPb{OrigName: true, EmitDefaults: true}
	yaml := &YamlMarshaler{}
	opt := func(mux *runtime.ServeMux) {
		runtime.WithMarshalerOption("application/json", json)(mux)
		runtime.WithMarshalerOption("application/yaml", yaml)(mux)
	}
	mux := runtime.NewServeMux(opt)

	// Configure the gateway.
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterLivingRoomHandlerFromEndpoint(ctx, mux, host+":"+grpcPort, opts); err != nil {
		logger.Error(err, "failed to start REST gateway")
	}

	// Return the handler.
	return mux
}

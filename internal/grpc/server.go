package grpc

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/sebastianrosch/couchconnections/internal/service"
	"github.com/twitchtv/twirp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/sebastianrosch/couchconnections/rpc/couchconnections-api/v1"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

// Authenticator interface
type Authenticator interface {
	// Authenticate authenticates a request
	Authenticate(ctx context.Context) (context.Context, error)
}

// GetServer returns the gRPC server and publishes the procedure endpoints.
func GetServer(
	ctx context.Context,
	logger logr.Logger,
	v1Service v1.CouchConnectionsServer,
	authenticator Authenticator) *grpc.Server {
	authenticatorMiddleware := authenticatorAsUnaryInterceptor(authenticator)

	// Register the gRPC server.
	middlewares := grpc_middleware.ChainUnaryServer(extractMethodInfoMiddleware, authenticatorMiddleware, convertTwirpError)
	server := grpc.NewServer(grpc.UnaryInterceptor(middlewares))
	v1.RegisterCouchConnectionsServer(server, v1Service)

	// Return the gRPC server.
	return server
}

// convertTwirpError converts the internal Twirp error to a gRPC error code, if one occurred.
func convertTwirpError(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)

	if err != nil {
		if twerr, ok := err.(twirp.Error); ok {
			switch twerr.Code() {
			case twirp.Canceled:
				return resp, status.Errorf(codes.Canceled, twerr.Msg())
			case twirp.Unknown:
				return resp, status.Errorf(codes.Unknown, twerr.Msg())
			case twirp.InvalidArgument:
				return resp, status.Errorf(codes.InvalidArgument, twerr.Msg())
			case twirp.DeadlineExceeded:
				return resp, status.Errorf(codes.DeadlineExceeded, twerr.Msg())
			case twirp.NotFound:
				return resp, status.Errorf(codes.NotFound, twerr.Msg())
			case twirp.AlreadyExists:
				return resp, status.Errorf(codes.AlreadyExists, twerr.Msg())
			case twirp.PermissionDenied:
				return resp, status.Errorf(codes.PermissionDenied, twerr.Msg())
			case twirp.ResourceExhausted:
				return resp, status.Errorf(codes.ResourceExhausted, twerr.Msg())
			case twirp.FailedPrecondition:
				return resp, status.Errorf(codes.FailedPrecondition, twerr.Msg())
			case twirp.Aborted:
				return resp, status.Errorf(codes.Aborted, twerr.Msg())
			case twirp.OutOfRange:
				return resp, status.Errorf(codes.OutOfRange, twerr.Msg())
			case twirp.Unimplemented:
				return resp, status.Errorf(codes.Unimplemented, twerr.Msg())
			case twirp.Internal:
				return resp, status.Errorf(codes.Internal, twerr.Msg())
			case twirp.Unavailable:
				return resp, status.Errorf(codes.Unavailable, twerr.Msg())
			case twirp.DataLoss:
				return resp, status.Errorf(codes.DataLoss, twerr.Msg())
			case twirp.Unauthenticated:
				return resp, status.Errorf(codes.Unauthenticated, twerr.Msg())
			}
		}
	}

	return resp, err
}

// extractMethodInfoMiddleware extracts the full method name and it stores into the context
func extractMethodInfoMiddleware(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	methodInfo := &service.MethodInfo{
		FullName: info.FullMethod,
	}
	ctx = service.NewMetadata().WithMethodInfo(ctx, methodInfo)
	return handler(ctx, req)
}

// authenticatorAsUnaryInterceptor calls the Authenticate function and wraps the error as a grpc Unauthenticated error
func authenticatorAsUnaryInterceptor(authenticator Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := authenticator.Authenticate(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}

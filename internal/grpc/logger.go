package grpc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

// RequestLogger logs grpc requests and responses.
type RequestLogger struct {
	logger      logr.Logger
	apiEndpoint string
}

// NewRequestLogger returns a new RequestLogger.
func NewRequestLogger(logger logr.Logger, apiEndpoint string) *RequestLogger {
	return &RequestLogger{
		logger:      logger,
		apiEndpoint: apiEndpoint,
	}
}

// LogGRPCCall logs the gRPC call.
func (l *RequestLogger) LogGRPCCall(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	l.logger.Info(fmt.Sprintf("%s%s", l.apiEndpoint, method))
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

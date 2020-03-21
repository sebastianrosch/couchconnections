package grpc

import (
	"context"
	"fmt"

	"github.com/twitchtv/twirp"
)

// GetMethodInfo returns the gRPC method info from the context.
func GetMethodInfo(ctx context.Context) string {
	pkg, ok := twirp.PackageName(ctx)
	if !ok {
		return ""
	}
	svc, ok := twirp.ServiceName(ctx)
	if !ok {
		return ""
	}
	method, ok := twirp.MethodName(ctx)
	if !ok {
		return ""
	}
	return fmt.Sprintf("/%s.%s/%s", pkg, svc, method)
}

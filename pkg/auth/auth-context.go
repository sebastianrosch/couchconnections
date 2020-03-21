package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// AuthContextError is the struct for authentication context errors
type AuthContextError struct { //nolint
	Message string
}

// Error returns the message of the AuthContextError
func (a AuthContextError) Error() string {
	return a.Message
}

// BearerTokenContext is the struct that handles adding and retrieving the bearer authorization header from context
type BearerTokenContext struct {
}

// GetAuthTokenFromAuthorizationHeader returns the auth token from the authorization header
func (b *BearerTokenContext) GetAuthTokenFromAuthorizationHeader(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok && len(md.Get("authorization")) != 0 {
		authHeader := strings.TrimPrefix(md.Get("authorization")[0], "Bearer ")

		return strings.TrimPrefix(authHeader, "bearer ")
	}

	return ""
}

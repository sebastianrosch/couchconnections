package auth

import (
	"context"
	"errors"

	"github.com/go-logr/logr"

	"github.com/sebastianrosch/livingroompresentations/internal/service"
)

// MetadataRetriever interface provides access to the metadata of the current request
type MetadataRetriever interface {
	// GetMethodInfo returns the method information from the context
	GetMethodInfo(ctx context.Context) *service.MethodInfo
}

// TokenContext provides a function to get the token from the context.
type TokenContext interface {
	GetAuthTokenFromAuthorizationHeader(ctx context.Context) string
}

// TokenAuthenticator implements token authentication
type TokenAuthenticator struct {
	logger            logr.Logger
	whitelist         []string
	tokenDecoder      *JWTTokenDecoder
	userInfoRetriever *UserInfoRetriever
	metadata          MetadataRetriever
	tokenContext      TokenContext
}

type userInfoKey struct{}

// Authenticate authenticates a request by validating the "authorization" header from the request metadata
func (t *TokenAuthenticator) Authenticate(ctx context.Context) (context.Context, error) {
	if t.skipAuthentication(ctx) {
		return ctx, nil
	}

	tokenString := t.tokenContext.GetAuthTokenFromAuthorizationHeader(ctx)

	userInfo, err := t.userInfoRetriever.GetUserInfo(ctx, tokenString)
	if err != nil {
		if uerr, ok := err.(UserInfoAuthorizationError); ok {
			// We're logging this as info because it is most likely a problem with the token.
			t.logger.Info("authentication failed while getting user info", "error", uerr)
			return nil, errors.New(("invalid auth token"))
		}
		return nil, err
	}

	accessTokenClaims, err := t.tokenDecoder.DecodeAndValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, userInfoKey{}, userInfo)
	ctx = WithAuthorizationPermissions(ctx, accessTokenClaims.Permissions)

	return ctx, nil
}

// skipAuthentication returns true if the current method is whitlisted from authentication
func (t *TokenAuthenticator) skipAuthentication(ctx context.Context) bool {
	methodInfo := t.metadata.GetMethodInfo(ctx)
	if methodInfo == nil {
		return false
	}

	for _, whitelisted := range t.whitelist {
		if whitelisted == methodInfo.FullName {
			return true
		}
	}
	return false
}

// NewAuthenticator returns a new Authenticator
func NewAuthenticator(
	logger logr.Logger,
	whitelist []string,
	tokenDecoder *JWTTokenDecoder,
	userInfoRetriever *UserInfoRetriever,
	metadata MetadataRetriever,
	tokenContext TokenContext) *TokenAuthenticator {
	return &TokenAuthenticator{
		logger:            logger,
		whitelist:         whitelist,
		tokenDecoder:      tokenDecoder,
		userInfoRetriever: userInfoRetriever,
		metadata:          metadata,
		tokenContext:      tokenContext,
	}
}

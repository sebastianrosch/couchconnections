package auth

import (
	"context"
	"io"
	"net/http"

	"github.com/sebastianrosch/livingroompresentations/pkg/auth/internal"
)

// Utils interface
type Utils interface {
	RandomBytes(length int) []byte
	Encode(msg []byte) string
	Sha256Hash(value string) string
	DecodeJSON(reader io.Reader, into interface{}) error
	PostForm(ctx context.Context, url string, data map[string][]string, httpClient *http.Client) (resp *http.Response, err error)
	GetURLWithHeaders(ctx context.Context, url string, headers map[string]string, httpClient *http.Client) (resp *http.Response, err error)
	GetRedirectURLFromConfig(config internal.RedirectURLProvider) string
	GetLogoutURLFromConfig(config internal.LogoutURLProvider) string
	ListenSingleRequest(ctx context.Context, address string, port string, endpoint string, handler http.HandlerFunc)
	DefaultHttpClient() *http.Client
}

// URLOpener is used to open URLs
type URLOpener interface {
	OpenURL(url string)
}

// HTMLRenderer interface
type HTMLRenderer interface {
	RenderSuccessPage(writer io.Writer, emailAddress, logoutURL string) error
	RenderErrorPage(writer io.Writer, errorMessage string) error
}

// AuthorizationCallbackChannel struct
type AuthorizationCallbackChannel struct {
	Token *TokenResponse
	Error error
}

// CallbackHandlerBuilder interface
type CallbackHandlerBuilder interface {
	BuildCallbackHandler(channel chan<- AuthorizationCallbackChannel, state, verifier, redirectURL, logoutURL string) http.HandlerFunc
}

// CodeVerifier interface
type CodeVerifier interface {
	CreateCodeChallenge(method string) (*internal.CodeChallenge, error)
}

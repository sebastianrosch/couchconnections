package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sebastianrosch/livingroompresentations/pkg/auth/internal"
)

type tokenEndpointResponse struct {
	TokenResponse
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// TokenRetrieverConfig contains the configurations for the token retriever
type TokenRetrieverConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenRetriever exposes different functions that can be used for retrieving tokens from the IDP
type TokenRetriever struct {
	config     *TokenRetrieverConfig
	utils      Utils
	httpClient *http.Client
}

// AccessCode Returns an [access token](#type-tokenresponse) for a given access code and code verifier (PKCE flow).
func (t *TokenRetriever) AccessCode(ctx context.Context, accessCode, codeVerifier, redirectURL string) (*TokenResponse, error) {
	formData := map[string][]string{
		"grant_type":    {"authorization_code"},
		"client_id":     {t.config.ClientID},
		"code_verifier": {codeVerifier},
		"code":          {accessCode},
		"redirect_uri":  {redirectURL},
	}

	return t.tokenEndpointRequest(ctx, formData)
}

// RefreshToken requests an [access token](#type-tokenresponse) for a given refresh token (refresh token flow).
func (t *TokenRetriever) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	formData := map[string][]string{
		"grant_type":    {"refresh_token"},
		"client_id":     {t.config.ClientID},
		"refresh_token": {refreshToken},
	}

	return t.tokenEndpointRequest(ctx, formData)
}

// TokenExchange requests a new [access token](#type-tokenresponse) to use with a target audience given an access token (token exhange flow).
func (t *TokenRetriever) TokenExchange(ctx context.Context, accessToken string, audience string, scope []string) (*TokenResponse, error) {
	formData := map[string][]string{
		"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"client_id":          {t.config.ClientID},
		"subject_token":      {accessToken},
		"subject_token_type": {"urn:ietf:params:oauth:token-type:access_token"},
		"audience":           {audience},
		"client_secret":      {t.config.ClientSecret},
		"scope":              {strings.Join(scope, " ")},
	}

	return t.tokenEndpointRequest(ctx, formData)
}

func (t *TokenRetriever) tokenEndpointRequest(ctx context.Context, formData map[string][]string) (*TokenResponse, error) {
	resp, err := t.utils.PostForm(ctx, t.config.TokenURL, formData, t.httpClient)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result tokenEndpointResponse
	t.utils.DecodeJSON(resp.Body, &result)

	if result.Error != "" {
		return nil, fmt.Errorf("%s: %s", result.Error, result.ErrorDescription)
	}

	return &TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		IDToken:      result.IDToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// NewTokenRetriever returns a new instance of [TokenRetriever](#type-tokenretriever)
func NewTokenRetriever(config *TokenRetrieverConfig, httpClient *http.Client) *TokenRetriever {
	utils := internal.NewUtils()

	return NewTokenRetrieverWithBoundaries(config, utils, httpClient)
}

// NewTokenRetrieverWithBoundaries rreturns a new instance of [TokenRetriever](#type-tokenretriever) with the provided boundaries
func NewTokenRetrieverWithBoundaries(config *TokenRetrieverConfig, utils Utils, httpClient *http.Client) *TokenRetriever {
	return &TokenRetriever{config: config, utils: utils, httpClient: httpClient}
}

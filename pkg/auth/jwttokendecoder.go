package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

// SigningKey struct
type SigningKey struct {
	Alg string
	Kty string
	Use string
	X5c []string
	N   string
	E   string
	Kid string
	X5t string
}

// JwksResponse response from the auth server for the JWKS endpoint
type JwksResponse struct {
	Keys []SigningKey
}

// AccessTokenClaims represents claims usually present on an access token
type AccessTokenClaims struct {
	Aud         []string
	Azp         string
	Exp         int
	Iat         int
	Iss         string
	Permissions []string `mapstructure:"http://livingroom/roles"`
	Scope       string
	Sub         string
}

// IDTokenClaims represents claims usually present on an IDToken
type IDTokenClaims struct {
	Iss           string
	Sub           string
	Aud           string
	Iat           int
	Exp           int
	Email         string
	EmailVerified bool   `mapstructure:"email_verified"`
	GivenName     string `mapstructure:"given_name"`
	FamilyName    string `mapstructure:"family_name"`
	Nickname      string
	Name          string
	Picture       string
	Locale        string
	UpdatedAt     string `mapstructure:"updated_at"`
}

type validSigningKey struct {
	kid       string
	publicKey string
}

// JWTTokenDecoder provides functions to validate and decode JWT
type JWTTokenDecoder struct {
	jwksURL string
}

// DecodeAndValidate validates the jwt and returns the token parsed
func (t *JWTTokenDecoder) DecodeAndValidate(tokenString string) (jwt.MapClaims, error) {
	token, err := t.extractAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}

	return nil, fmt.Errorf("Error decoding claims")
}

// DecodeAndValidateAccessToken validates the jwt and returns the token parsed as [access token](#type-accesstokenclaims)
func (t *JWTTokenDecoder) DecodeAndValidateAccessToken(accessTokenString string) (*AccessTokenClaims, error) {
	claims, err := t.DecodeAndValidate(accessTokenString)
	if err != nil {
		return nil, err
	}

	var accessTokenClaims AccessTokenClaims
	err = mapstructure.Decode(claims, &accessTokenClaims)
	if err != nil {
		return nil, err
	}

	return &accessTokenClaims, nil
}

// DecodeAndValidateIDToken validates the jwt and returns the token parsed as id token
func (t *JWTTokenDecoder) DecodeAndValidateIDToken(idTokenString string) (*IDTokenClaims, error) {
	claims, err := t.DecodeAndValidate(idTokenString)
	if err != nil {
		return nil, err
	}

	var idTokenClaims IDTokenClaims
	err = mapstructure.Decode(claims, &idTokenClaims)
	if err != nil {
		return nil, err
	}

	return &idTokenClaims, nil
}
func (t *JWTTokenDecoder) extractAndValidateToken(tokenString string) (*jwt.Token, error) {
	// Add a 60 second leeway to prevent possible clock skew issues
	jwt.TimeFunc = func() time.Time {
		leeway := time.Minute
		return time.Now().Add(leeway)
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)

		if !ok {
			return nil, fmt.Errorf("Kid not present in token headers")
		}

		pem, err := t.getSigningKey(kid)

		if err != nil {
			return nil, fmt.Errorf("Error getting signing key: %s", err)
		}

		return jwt.ParseRSAPublicKeyFromPEM([]byte(pem.publicKey))
	})
}

func (t *JWTTokenDecoder) getSigningKey(kid string) (*validSigningKey, error) {
	keys, err := t.getSigningKeys()

	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		if k.kid == kid {
			return &k, nil
		}
	}

	return nil, fmt.Errorf("Unable to find a signing key that matches %s", kid)
}

func (t *JWTTokenDecoder) getSigningKeys() ([]validSigningKey, error) {
	keys, err := t.getJWKS()

	if err != nil {
		return nil, err
	}

	var filteredSigningKeys []validSigningKey

	for _, k := range keys {
		if k.Use == "sig" &&
			k.Kty == "RSA" &&
			k.Kid != "" &&
			((k.X5c != nil && len(k.X5c) > 0) || (k.N != "" && k.E != "")) {
			validKey := validSigningKey{
				kid:       k.Kid,
				publicKey: t.certToPEM(k.X5c[0]),
			}
			filteredSigningKeys = append(filteredSigningKeys, validKey)
		}
	}

	return filteredSigningKeys, nil
}

func (t *JWTTokenDecoder) getJWKS() ([]SigningKey, error) {
	response, err := http.Get(t.jwksURL)

	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("Wrong response. Status code: %d", response.StatusCode)
	}

	var jwks JwksResponse

	json.NewDecoder(response.Body).Decode(&jwks)

	return jwks.Keys, nil
}

func (t *JWTTokenDecoder) certToPEM(cert string) string {
	re := regexp.MustCompile(`.{1,64}`)
	matches := re.FindAll([]byte(cert), -1)
	c := bytes.Join(matches, []byte("\n"))

	return fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----\n", c)
}

// NewJWTTokenDecoder returns a default token parser
func NewJWTTokenDecoder(jwksURL string) *JWTTokenDecoder {
	return &JWTTokenDecoder{jwksURL: jwksURL}
}

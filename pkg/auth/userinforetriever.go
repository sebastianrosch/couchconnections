package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sebastianrosch/couchconnections/pkg/auth/internal"
)

// UserInfoAuthorizationError represents an authorization error that occurred while
// getting the user info from the identity provider
type UserInfoAuthorizationError struct {
	message string
}

// Error returns the error message
func (e UserInfoAuthorizationError) Error() string {
	return e.message
}

// NewUserInfoAuthorizationError returns a new instance of UserInfoAuthorizationError
func NewUserInfoAuthorizationError(message string) *UserInfoAuthorizationError {
	return &UserInfoAuthorizationError{message: message}
}

// UserInfoResponse struct with the user information returned by [GetUserInfo](#func-getuserinfo)
type UserInfoResponse struct {
	Sub                 string
	Name                string
	GivenName           string `json:"given_name"`
	FamilyName          string `json:"family_name"`
	MiddleName          string `json:"middle_name"`
	Nickname            string
	PreferredUsername   string `json:"preferred_username"`
	Profile             string
	Picture             string
	Website             string
	Email               string
	EmailVerified       bool `json:"email_verified"`
	Gender              string
	Birthdate           string
	Zoneinfo            string
	Locale              string
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Address             UserInfoAddress
	UpdatedAt           string `json:"updated_at"`
}

// UserInfoAddress contains the address information of the user
type UserInfoAddress struct {
	Country string
}

// UserInfoRetriever struct exposes a func to retrieve the user info
type UserInfoRetriever struct {
	userInfoEndpoint string
	utils            Utils
	httpClient       *http.Client
}

// GetUserInfo queries the IDP for [user information](#type-userinforesponse) of the provided access token. This function can also be used to validate an opaque access token.
func (u *UserInfoRetriever) GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponse, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	res, err := u.utils.GetURLWithHeaders(ctx, u.userInfoEndpoint, headers, u.httpClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// http does not throw an error on authorization failures, so we need to check the status code.
	if res.StatusCode >= 200 && res.StatusCode <= 299 {

		var userInfo UserInfoResponse

		err = u.utils.DecodeJSON(res.Body, &userInfo)

		if err != nil {
			return nil, err
		}

		return &userInfo, nil
	} else if res.StatusCode >= 400 && res.StatusCode <= 499 {
		return nil, NewUserInfoAuthorizationError("An authorization error occurred while getting the user info from the identity provider. Please make sure that the token is valid.")
	} else {
		return nil, fmt.Errorf("an unknown error occurred while getting the user info from the identity provider: status code %d", res.StatusCode)
	}
}

// NewUserInfoRetriever returns a new instance of [UserInfoRetriever](#type-userinforetriever)
func NewUserInfoRetriever(userInfoEndpoint string, httpClient *http.Client) *UserInfoRetriever {
	utils := internal.NewUtils()

	return NewUserInfoRetrieverWithBoundaries(userInfoEndpoint, utils, httpClient)
}

// NewUserInfoRetrieverWithBoundaries returns a new instance of [UserInfoRetriever](#type-userinforetriever) with the provided boundaries
func NewUserInfoRetrieverWithBoundaries(userInfoEndpoint string, utils Utils, httpClient *http.Client) *UserInfoRetriever {
	return &UserInfoRetriever{userInfoEndpoint: userInfoEndpoint, utils: utils, httpClient: httpClient}
}

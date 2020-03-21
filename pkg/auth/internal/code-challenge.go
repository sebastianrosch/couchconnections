package internal

import (
	"fmt"
)

// CodeChallenge struct
type CodeChallenge struct {
	Verifier  string
	Challenge string
	Method    string
}

// CodeVerifier struct provides PKCE code verifier operations
type CodeVerifier struct {
	utils CodeChallengeUtils
}

// NewCodeVerifier func returns a new instance of CodeVerifier
func NewCodeVerifier(utils CodeChallengeUtils) *CodeVerifier {
	return &CodeVerifier{
		utils: utils,
	}
}

// CreateCodeChallenge func creates a CodeChallenge for the CodeVerifier
func (v *CodeVerifier) CreateCodeChallenge(method string) (*CodeChallenge, error) {
	if method == "plain" {
		return v.generateCodeChallengePlain(), nil
	}

	if method == "S256" {
		return v.generateCodeChallengeS256(), nil
	}

	return nil, fmt.Errorf("invalid method: %s", method)
}

func (v *CodeVerifier) generateCodeChallengePlain() *CodeChallenge {
	verifier := v.utils.Encode(v.utils.RandomBytes(32))
	return &CodeChallenge{
		Verifier:  verifier,
		Challenge: verifier,
		Method:    "plain",
	}
}

func (v *CodeVerifier) generateCodeChallengeS256() *CodeChallenge {
	verifier := v.utils.Encode(v.utils.RandomBytes(32))
	return &CodeChallenge{
		Verifier:  verifier,
		Challenge: v.utils.Sha256Hash(verifier),
		Method:    "S256",
	}
}

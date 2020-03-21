package internal

//go:generate mockgen -source=internal.go -destination=internal_mock.go -package=internal
type CodeChallengeUtils interface {
	RandomBytes(length int) []byte
	Encode(msg []byte) string
	Sha256Hash(value string) string
}

type LogoutURLProvider interface {
	GetLogoutURLParts() (url string, clientID string, returnURL string)
}

type RedirectURLProvider interface {
	GetRedirectURLParts() (address string, port string, endpoint string)
}

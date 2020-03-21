package internal

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// RandGenerator interface
type RandGenerator interface {
	Intn(number int) int
}

type Utils struct {
}

// RandomBytes func
func (u *Utils) RandomBytes(length int) []byte {
	generator := u.newRand()
	bytes := make([]byte, length, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(generator.Intn(255))
	}

	return bytes
}

// Encode func
func (u *Utils) Encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)

	return encoded
}

// Sha256Hash func
func (u *Utils) Sha256Hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	return u.Encode(hash.Sum(nil))
}

func (u *Utils) DecodeJSON(reader io.Reader, into interface{}) error {
	return json.NewDecoder(reader).Decode(into)
}

func (u *Utils) PostForm(ctx context.Context, url string, data map[string][]string, httpClient *http.Client) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, encodeFormBody(data))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if httpClient == nil {
		httpClient = u.DefaultHttpClient()
	}
	return httpClient.Do(req)
}

func (u *Utils) GetURLWithHeaders(ctx context.Context, url string, headers map[string]string, httpClient *http.Client) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if httpClient == nil {
		httpClient = u.DefaultHttpClient()
	}
	return httpClient.Do(req)
}

func (u *Utils) GetRedirectURLFromConfig(config RedirectURLProvider) string {
	address, port, endpoint := config.GetRedirectURLParts()
	return fmt.Sprintf("http://%s:%s%s", address, port, endpoint)
}

func (u *Utils) GetLogoutURLFromConfig(config LogoutURLProvider) string {
	logoutURL, clientID, returnURL := config.GetLogoutURLParts()
	return fmt.Sprintf("%s?client_id=%s&returnTo=%s", logoutURL, url.QueryEscape(clientID), url.QueryEscape(returnURL))
}

func (u *Utils) ListenSingleRequest(ctx context.Context, address string, port string, endpoint string, handler http.HandlerFunc) {
	// Configure the listener
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Println("ERROR", err)
		panic(err)
	}

	// Configure the handler function
	router := http.NewServeMux()
	router.HandleFunc(endpoint, handler)

	// Configure the server
	server := &http.Server{
		Handler: router,
	}

	// Start the server
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			fmt.Println("ERROR", err)
			panic(err)
		}
	}()

	// Close the server in case the context gets canceled
	go func() {
		<-ctx.Done()
		server.Close()
	}()
}

func (u *Utils) DefaultHttpClient() *http.Client {
	// extracted from https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: defaultTransport(),
	}
}

func (u *Utils) newRand() RandGenerator {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// defaultTransport returns a new http.Transport with similar default values to
// http.DefaultTransport, but with idle connections and keepalives disabled.
func defaultTransport() *http.Transport {
	transport := defaultPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// defaultPooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func defaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}

func encodeFormBody(data url.Values) io.Reader {
	return strings.NewReader(data.Encode())
}

// NewUtils func
func NewUtils() *Utils {
	return &Utils{}
}

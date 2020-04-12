package opentok

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// OpenTok API host URL
const defaultAPIHost = "https://api.opentok.com"

// For use in X-TB-TOKEN-AUTH header value
const tokenSentinel = "T1=="

type issueType string

type HttpDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

const (
	// For most REST API calls, set issue type to "project"
	projectToken issueType = "project"
	// For Account Management REST methods, set issue type to "account"
	accountToken issueType = "account"
)

// OpenTok stores the API key and secret for use in making API call
type OpenTok struct {
	apiKey    string
	apiSecret string
	apiHost   string

	httpClient HttpDoer
}

// New returns an initialized OpenTok instance with the API key and API secret.
func New(apiKey, apiSecret string) *OpenTok {
	return &OpenTok{apiKey, apiSecret, defaultAPIHost, http.DefaultClient}
}

// SetAPIHost is used to set OpenTok API Host to specific URL
func (ot *OpenTok) SetAPIHost(url string) error {
	if url == "" {
		return fmt.Errorf("OpenTok API Host cannot be empty")
	}

	ot.apiHost = url

	return nil
}

// SetHttpClient specifies http client, http.DefaultClient used by default.
func (ot *OpenTok) SetHttpClient(client HttpDoer) {
	if client != nil {
		ot.httpClient = client
	}
}

// Generate JWT token for API calls
func (ot *OpenTok) jwtToken(ist issueType) (string, error) {
	type OpenTokClaims struct {
		Ist issueType `json:"ist,omitempty"`
		jwt.StandardClaims
	}

	issuedAt := time.Now().UTC()

	claims := OpenTokClaims{
		ist,
		jwt.StandardClaims{
			Issuer:    ot.apiKey,
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: issuedAt.Add((5 * time.Minute)).Unix(), // The maximum allowed expiration time range is 5 minutes.
			Id:        uuid.New().String(),
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the api secret
	return token.SignedString([]byte(ot.apiSecret))
}

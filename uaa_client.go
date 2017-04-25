package autoscaler

import (
	"fmt"
	"net/http"

	"crypto/tls"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Endpoint to Cloud Controller
type Endpoint struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

//AccessToken - represents a authenticated token from UAA
type AccessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string
	JTI       string
}

// CFConfig Config represents all the configuration for making a oauth2 call
type CFConfig struct {
	CCApiURL          string
	Username          string
	Password          string
	ClientID          string
	ClientSecret      string
	SkipSslValidation bool
	httpClient        *http.Client
	TokenSource       oauth2.TokenSource
}

// OauthHTTPWrapper is an http client wrapper that makes the call with an oauth2 token
type OauthHTTPWrapper interface {
	NewCCRequest(method, path string, body io.Reader) (*http.Request, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	Do(request *http.Request) (*http.Response, error)
}

// TokenHandlingClient - UAA Client
type TokenHandlingClient struct {
	Config   *CFConfig
	Endpoint *Endpoint
}

// DefaultCFConfig - default configuraiton for making CF calls
func DefaultCFConfig() *CFConfig {
	return &CFConfig{
		CCApiURL:          "https://api.local.pcfdev.io",
		Username:          "admin",
		Password:          "admin",
		SkipSslValidation: true,
		httpClient:        http.DefaultClient,
	}
}

// NewUAAClient - Creates a new UAA Client
func NewUAAClient(config *CFConfig) (OauthHTTPWrapper, error) {
	ctx := context.Background()
	defConfig := DefaultCFConfig()

	if !config.SkipSslValidation {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, defConfig.httpClient)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: tr})
	}

	endpoint, err := getInfo(config.CCApiURL, oauth2.NewClient(ctx, nil))

	if err != nil {
		return nil, fmt.Errorf("Could not get api /v2/info: %v", err)
	}

	switch {
	case config.ClientID != "":
		config = getClientAuth(config, endpoint, ctx)
	default:
		config, err = getUserAuth(config, endpoint, ctx)
		if err != nil {
			return nil, err
		}
	}

	client := &TokenHandlingClient{
		Config:   config,
		Endpoint: endpoint,
	}
	return client, nil
}

// Do ...
func (client *TokenHandlingClient) Do(request *http.Request) (*http.Response, error) {
	return client.Config.httpClient.Do(request)
}

// NewCCRequest ...
func (client *TokenHandlingClient) NewCCRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", client.Config.CCApiURL, path), body)
}

// NewRequest ...
func (client *TokenHandlingClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func getUserAuth(config *CFConfig, endpoint *Endpoint, ctx context.Context) (*CFConfig, error) {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  endpoint.AuthorizationEndpoint + "/oauth/auth",
			TokenURL: endpoint.TokenEndpoint + "/oauth/token",
		},
	}

	token, err := authConfig.PasswordCredentialsToken(ctx, config.Username, config.Password)

	if err != nil {
		return nil, fmt.Errorf("Error getting token: %v", err)
	}

	config.TokenSource = authConfig.TokenSource(ctx, token)
	config.httpClient = oauth2.NewClient(ctx, config.TokenSource)

	return config, err
}

func getClientAuth(config *CFConfig, endpoint *Endpoint, ctx context.Context) *CFConfig {
	authConfig := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     endpoint.TokenEndpoint + "/oauth/token",
	}

	config.TokenSource = authConfig.TokenSource(ctx)
	config.httpClient = authConfig.Client(ctx)
	return config
}

func getInfo(api string, httpClient *http.Client) (*Endpoint, error) {
	var endpoint Endpoint

	if api == "" {
		return nil, errors.New("Missing CC API url")
	}

	resp, err := httpClient.Get(api + "/v2/info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, err
}

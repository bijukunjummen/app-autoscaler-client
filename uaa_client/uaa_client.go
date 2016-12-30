package uaa_client

import (
	"fmt"
	"net/http"

	"crypto/tls"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
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

// Config represents all the configuration for making a oauth2 call
type Config struct {
	CCApiUrl          string
	Username          string
	Password          string
	ClientId          string
	ClientSecret      string
	SkipSslValidation bool
	httpClient        *http.Client
	TokenSource       oauth2.TokenSource
}

type OauthHttpWrapper interface {
	NewCCRequest(method, path string, body io.Reader) (*http.Request, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	Do(request *http.Request) (*http.Response, error)
}

// Client - UAA Client
type TokenHandlingClient struct {
	Config   *Config
	Endpoint *Endpoint
}

func DefaultConfig() *Config {
	return &Config{
		CCApiUrl:          "https://api.local.pcfdev.io",
		Username:          "admin",
		Password:          "admin",
		SkipSslValidation: true,
		httpClient:        http.DefaultClient,
	}
}

// NewClient - Creates a new UAA Client
func NewClient(config *Config) (OauthHttpWrapper, error) {
	ctx := context.Background()
	defConfig := DefaultConfig()

	if !config.SkipSslValidation {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, defConfig.httpClient)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: tr})
	}

	endpoint, err := getInfo(config.CCApiUrl, oauth2.NewClient(ctx, nil))

	if err != nil {
		return nil, fmt.Errorf("Could not get api /v2/info: %v", err)
	}

	switch {
	case config.ClientId != "":
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

func (client *TokenHandlingClient) Do(request *http.Request) (*http.Response, error) {
	return client.Config.httpClient.Do(request)
}

func (client *TokenHandlingClient) NewCCRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", client.Config.CCApiUrl, path), body)
}

func (client *TokenHandlingClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func getUserAuth(config *Config, endpoint *Endpoint, ctx context.Context) (*Config, error) {
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

func getClientAuth(config *Config, endpoint *Endpoint, ctx context.Context) *Config {
	authConfig := &clientcredentials.Config{
		ClientID:     config.ClientId,
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

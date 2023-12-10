package bunce_go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BunceBaseURL = "https://api.bunce.so"
	Version      = "v1"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

type Config struct {
	baseURL    *url.URL
	HttpClient *http.Client
}

type Client struct {
	apiKey string
	config *Config
	common service

	Customers *Customer
}

type service struct {
	client *Client
}

func New(apiKey string, config *Config) *Client {
	config.baseURL = buildBaseUrl(config)

	if config.HttpClient == nil {
		config.HttpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	c := &Client{
		apiKey: apiKey,
		config: config,
	}
	c.common.client = c
	c.Customers = newCustomer(c)

	return c
}

func (c Client) decode(v interface{}, b []byte) (err error) {
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

func buildBaseUrl(config *Config) *url.URL {
	if config.baseURL == nil {
		rawURL := fmt.Sprintf("%s/%s", BunceBaseURL, Version)
		return MustParseURL(rawURL)
	}

	if strings.Contains(config.baseURL.String(), "bunce.so/v") {
		return config.baseURL
	}

	return config.baseURL.JoinPath(Version)
}

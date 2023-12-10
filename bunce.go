package bunce_go

import (
	"encoding/json"
	"net/http"
	"net/url"
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

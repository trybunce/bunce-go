package bunce_go

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func (c Client) sendRequest(req *http.Request, result interface{}) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Authorization", fmt.Sprintf("%s", c.apiKey))

	res, err := c.config.HttpClient.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "failed to send request")
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return res, errors.Wrap(err, "error while reading the response bytes")
	}

	var response APIResponse

	err = c.decode(&response, body)
	if err != nil {
		return res, errors.Wrap(err, "unable to unmarshal response body")
	}

	if !response.Status && invalidStatusCode(res.StatusCode) {
		return res, fmt.Errorf("%s", response.Message)
	}

	if result != nil {
		err = json.Unmarshal(*response.Data, result)
		if err != nil {
			return res, fmt.Errorf("error while unmarshalling the response data bytes %+v ", err)
		}
	}

	return res, nil
}

func (c Client) decode(v interface{}, b []byte) (err error) {
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func invalidStatusCode(actual int) bool {
	expected := map[int]bool{
		200: true,
		202: true,
		204: true,
	}

	if _, ok := expected[actual]; ok {
		return false
	}

	return true
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

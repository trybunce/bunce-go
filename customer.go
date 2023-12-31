package bunce_go

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNo   string `json:"phone_no,omitempty"`
}

type UpdateCustomerRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	PhoneNo   *string `json:"phone_no,omitempty"`
}

type CreateCustomerResponse struct {
	ID                string     `json:"id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	PhoneNo           string     `json:"phone_no"`
	CustomerCreatedAt *time.Time `json:"customer_created_at"`
}

type BulkCreateCustomerRequest []CreateCustomerRequest

type BulkCreateCustomerResponsePayload []CreateCustomerResponse

type CustomerPayload struct {
	ID                string     `json:"id"`
	FirstName         *string    `json:"first_name"`
	LastName          *string    `json:"last_name"`
	Email             string     `json:"email"`
	PhoneNo           *string    `json:"phone_no"`
	Providers         *string    `json:"providers"`
	CustomerCreatedAt *time.Time `json:"customer_created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	LastInteraction   *time.Time `json:"last_interaction"`
}

type CustomersResponsePayload struct {
	Data []CustomerPayload `json:"data"`
	Meta Pagination        `json:"meta"`
}

type CompanyQueryOptions struct {
	Page     int    `queryKey:"page"`
	PerPage  int    `queryKey:"per_page"`
	Query    string `queryKey:"query"`
	Provider string `queryKey:"provider"`
	Emails   string `queryKey:"email"`
}

type Customer service

func newCustomer(client *Client) *Customer {
	return &Customer{
		client: client,
	}
}

func (c *Customer) Create(ctx context.Context, data CreateCustomerRequest) (CreateCustomerResponse, error) {
	URL := c.client.config.baseURL.JoinPath("customers")
	var resp CreateCustomerResponse

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))

	if err != nil {
		return resp, err
	}

	_, err = c.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Customer) BulkCreate(ctx context.Context, customers BulkCreateCustomerRequest) (BulkCreateCustomerResponsePayload, error) {
	URL := c.client.config.baseURL.JoinPath("customers", "bulk")
	var resp BulkCreateCustomerResponsePayload

	jsonBody, err := json.Marshal(customers)
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = c.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Customer) Update(ctx context.Context, email string, data UpdateCustomerRequest) (interface{}, error) {
	URL := c.client.config.baseURL.JoinPath("customers", email)
	var resp interface{}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL.String(), bytes.NewBuffer(jsonBody))

	if err != nil {
		return resp, err
	}

	_, err = c.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Customer) Find(ctx context.Context, email string) (CustomerPayload, error) {
	URL := c.client.config.baseURL.JoinPath("customers", email)
	var resp CustomerPayload

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)

	if err != nil {
		return resp, err
	}

	_, err = c.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Customer) All(ctx context.Context, opts *CompanyQueryOptions) (*CustomersResponsePayload, error) {
	URL := c.client.config.baseURL.JoinPath("customers")
	var resp CustomersResponsePayload

	if opts != nil {
		queryValues := URL.Query()
		params, err := GenerateQueryParamsFromStruct(*opts)

		if err != nil {
			return nil, err
		}

		for _, param := range params {
			queryValues.Add(param.Key, param.Value)
		}

		URL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)

	if err != nil {
		return nil, err
	}

	_, err = c.client.sendRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

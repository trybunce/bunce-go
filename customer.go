package bunce_go

import (
	"context"
	"net/http"
	"time"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNo   string `json:"phone_no,omitempty"`
}

type BulkCreateCustomerRequest []CreateCustomerRequest

type CustomerResponse struct {
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
	Data []CustomerResponse `json:"data"`
	Meta Pagination         `json:"meta"`
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

func (c *Customer) Create(ctx context.Context) {
	URL := c.client.config.baseURL.JoinPath("customers")
}

func (c *Customer) BulkCreate(ctx context.Context, customers BulkCreateCustomerRequest) (interface{}, error) {
	URL := c.client.config.baseURL.JoinPath("customers", "bulk")
}

func (c *Customer) Get(ctx context.Context, opts *CompanyQueryOptions) (*CustomersResponsePayload, error) {
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

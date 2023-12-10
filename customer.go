package bunce_go

type CreateCustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNo   string `json:"phone_no,omitempty"`
}

type Customer service

func newCustomer(client *Client) *Customer {
	return &Customer{
		client: client,
	}
}

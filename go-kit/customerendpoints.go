package customersvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

// Endpoints are the actions we would like to offer to the outside world.
// Because we implement a RESTful API on top of a CRUD service, our endpoints
// look similar to our service. This does not always need to be the case.

// Endpoints take a request, call service method(s), and return a result.

// Note that creating client endpoints is out-of-scope for this sample.

// CustomerEndpoints is a collection of all endpoints that we offer
type CustomerEndpoints struct {
	GetCustomersEndpoint   endpoint.Endpoint
	GetCustomerEndpoint    endpoint.Endpoint
	AddCustomerEndpoint    endpoint.Endpoint
	DeleteCustomerEndpoint endpoint.Endpoint
	PatchCustomerEndpoint  endpoint.Endpoint
}

// MakeCustomerServerEndpoints creates endpoints for a given service
func MakeCustomerServerEndpoints(s CustomerService) CustomerEndpoints {
	return CustomerEndpoints{
		GetCustomersEndpoint:   MakeGetCustomersEndpoint(s),
		GetCustomerEndpoint:    MakeGetCustomerEndpoint(s),
		AddCustomerEndpoint:    MakeAddCustomerEndpoint(s),
		DeleteCustomerEndpoint: MakeDeleteCustomerEndpoint(s),
		PatchCustomerEndpoint:  MakePatchCustomerEndpoint(s),
	}
}

// MakeGetCustomersEndpoint creates an endpoint for the "get customers" operation
func MakeGetCustomersEndpoint(s CustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCustomersRequest)
		c, e := s.GetCustomers(ctx, req.OrderBy)
		return getCustomersResponse{Customers: c, Err: e}, nil
	}
}

// MakeGetCustomerEndpoint creates an endpoint for the "get customer" operation
func MakeGetCustomerEndpoint(s CustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(customerIDRequest)
		c, e := s.GetCustomer(ctx, req.ID)
		return customerResponse{Customer: c, Err: e}, nil
	}
}

// MakeAddCustomerEndpoint creates an endpoint for the "add customer" operation
func MakeAddCustomerEndpoint(s CustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(customerRequest)
		c, e := s.AddCustomer(ctx, req.Customer)
		return customerResponse{Customer: c, Err: e}, nil
	}
}

// MakeDeleteCustomerEndpoint creates an endpoint for the "delete customer" operation
func MakeDeleteCustomerEndpoint(s CustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(customerIDRequest)
		e := s.DeleteCustomer(ctx, req.ID)
		return deleteCustomerResponse{Err: e}, nil
	}
}

// MakePatchCustomerEndpoint creates an endpoint for the "patch customer" operation
func MakePatchCustomerEndpoint(s CustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchCustomerRequest)
		c, e := s.PatchCustomer(ctx, req.CustomerID, req.Customer)
		return customerResponse{Customer: c, Err: e}, nil
	}
}

// Request and response types

type getCustomersRequest struct {
	OrderBy string
}

type getCustomersResponse struct {
	Customers []Customer `json:"customers,omitempty"`
	Err       error      `json:"err,omitempty"`
}

type patchCustomerRequest struct {
	CustomerID uuid.UUID
	Customer Customer
}

// For the sake of brevity, we combine requests and responses with identical
// input/output parameters. Typically, every operation receives/responds with
// a unique set of parameters. Therefore, each endpoint often has its own
// request and response struct.

type customerIDRequest struct {
	ID uuid.UUID
}

type customerResponse struct {
	Customer Customer `json:"customer,omitempty"`
	Err      error    `json:"err,omitempty"`
}

type customerRequest struct {
	Customer Customer
}

type deleteCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

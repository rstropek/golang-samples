package customersvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting indicates that there is an inconsistency between routs and handlers
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// Transports bind our endpoints to a concrete transport protocol like HTTP or gRPC. A single
// microservice can support multiple transports. In our case, our microservice just offers
// a single transport: HTTP

// MakeCustomerHTTPHandler creates a http.Handler for a given service
func MakeCustomerHTTPHandler(s CustomerService, logger log.Logger) http.Handler {
	// In this sample we use the Gorilla multiplexer
	r := mux.NewRouter()

	// Create endpoints for the given service
	e := MakeCustomerServerEndpoints(s)

	// Some server options...
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// Bind endpoints to URLs and HTTP methods

	r.Methods("GET").Path("/customers").Handler(httptransport.NewServer(
		e.GetCustomersEndpoint,
		decodeGetCustomersRequest,
		encodeCustomersResponse,
		options...,
	))
	r.Methods("GET").Path("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}").Handler(httptransport.NewServer(
		e.GetCustomerEndpoint,
		decodeGetCustomerRequest,
		encodeCustomerResponse,
		options...,
	))
	r.Methods("POST").Path("/customers").Handler(httptransport.NewServer(
		e.AddCustomerEndpoint,
		decodeAddCustomerRequest,
		encodeAddCustomerResponse,
		options...,
	))
	r.Methods("DELETE").Path("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}").Handler(httptransport.NewServer(
		e.DeleteCustomerEndpoint,
		decodeDeleteCustomerRequest,
		encodeDeleteCustomerResponse,
		options...,
	))
	r.Methods("PATCH").Path("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}").Handler(httptransport.NewServer(
		e.PatchCustomerEndpoint,
		decodePatchCustomerRequest,
		encodeCustomerResponse,
		options...,
	))

	return r
}

// Methods for translating HTTP requests/responses into/from endpoint requests/responses

func decodeGetCustomersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getCustomersRequest{OrderBy: r.FormValue("orderBy")}, nil
}

func decodeGetCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrBadRouting
	}
	return customerIDRequest{ID: cid}, nil
}

func decodeAddCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req customerRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Customer); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeDeleteCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrBadRouting
	}
	return customerIDRequest{ID: cid}, nil
}

func decodePatchCustomerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrBadRouting
	}
	req := patchCustomerRequest{CustomerID: cid}
	if e := json.NewDecoder(r.Body).Decode(&req.Customer); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeGetCustomersRequest(ctx context.Context, req *http.Request, request interface{}) error {
	return encodeRequest(ctx, req, request)
}

func encodeCustomerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/customers/"
	return encodeRequest(ctx, req, request)
}

func encodeCustomerIDRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(customerIDRequest)
	req.URL.Path = "/customers/" + r.ID.String()
	return encodeRequest(ctx, req, request)
}

func encodeCustomersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(getCustomersResponse)
	if ok && e.Err != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Customers)
}

func encodeAddCustomerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(customerResponse)
	if ok && e.Err != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Err, w)
		return nil
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/customers/%s", e.Customer.CustomerID.String()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Customer)
}

func encodeCustomerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(customerResponse)
	if ok && e.Err != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Err, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Customer)
}

type errorer interface {
	error() error
}

func encodeDeleteCustomerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(deleteCustomerResponse); ok && e.Err != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Err, w)
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	if _, ok := err.(ErrMissingMandatoryValue); ok {
		return http.StatusBadRequest
	}

	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidOrderBy, ErrInvalidHourlyRate, ErrGivenCustomerID, ErrInvalidCountry:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

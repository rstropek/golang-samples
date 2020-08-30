# Microservice with GoKit

## Initialize Folder

* Create new folder *go-kit*

* Initialize Go modules and get necessary packages

```bash
go mod init <your-identifier>

go get github.com/go-kit/kit
go get github.com/gorilla/mux
go get github.com/google/uuid
go get github.com/shopspring/decimal
```

## Add Business Logic

* Create *customerservice.go*

```go
package customersvc

import (
    "context"
    "errors"
    "fmt"
    "sort"
    "sync"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// This file contains the business logic of our microservice. It is just the logic 
// without any respect to transport-related (e.g. http, gRPC) topics.

var (
    // ErrGivenCustomerID indicates that a customer ID was specified where it is not allowed
    ErrGivenCustomerID = errors.New("customer ID must be empty")

    // ErrInvalidCountry indicates that the country code is not valid
    ErrInvalidCountry = errors.New("Country name must be three characters long (use ISO 3166-1 Alpha-3 code)")

    // ErrInvalidHourlyRate indicates that the hourly rate is not valid
    ErrInvalidHourlyRate = errors.New("Hourly rate must be >= 0")

    // ErrNotFound indicates that no customer with the given ID exists
    ErrNotFound = errors.New("not found")

    // ErrInvalidOrderBy indicates that the order-by clause is invalid
    ErrInvalidOrderBy = errors.New("Currently, we can only order by companyName")
)

// ErrMissingMandatoryValue indicates that a mandatory field does not have a value.
// This is a custom error type (for more see https://blog.golang.org/go1.13-errors)
type ErrMissingMandatoryValue struct {
    Field string
}

func (r ErrMissingMandatoryValue) Error() string {
    return fmt.Sprintf("Missing mandatory value (%s)", r.Field)
}

// Customer holds data of a customer record
type Customer struct {
    CustomerID  uuid.UUID       `json:"customerID,omitempty"`
    CompanyName string          `json:"customerName"`
    ContactName string          `json:"contactName"`
    Country     string          `json:"country"`
    HourlyRate  decimal.Decimal `json:"hourlyRate"`
}

// Note that all methods of the CustomerService receive a context. For more
// information about context in Go in general see https://golang.org/pkg/context/.
// For more about context and GoKit see https://gokit.io/examples/stringsvc.html#threading-a-context.

// CustomerService is a simple CRUD interface for customer management
type CustomerService interface {
    // GetCustomers returns a list of all customers.
    // The resulting list can optionally be sorted (orderBy parameter).
    GetCustomers(ctx context.Context, orderBy string) ([]Customer, error)

    // GetCustomer returns the customer with the given ID.
    GetCustomer(ctx context.Context, cid uuid.UUID) (Customer, error)

    // AddCustomer adds a new customer to the list.
    AddCustomer(ctx context.Context, c Customer) (Customer, error)

    // DeleteCustomer deletes the customer with the given ID.
    DeleteCustomer(ctx context.Context, cid uuid.UUID) error

    // PatchCustomers updates field in the customer with the given ID.
    PatchCustomer(ctx context.Context, cid uuid.UUID, c Customer) (Customer, error)
}

// customerRepository is an in-memory implementation of a customer repository.
type customerRepository struct {
    customers      map[uuid.UUID]Customer
    customersMutex *sync.RWMutex
}

// NewCustomerRepository creates a customer repository
func NewCustomerRepository() CustomerService {
    return customerRepository{
        customers:      make(map[uuid.UUID]Customer, 0),
        customersMutex: &sync.RWMutex{},
    }
}

func (s customerRepository) GetCustomers(ctx context.Context, orderBy string) ([]Customer, error) {
    // Lock customers while accessing it
    s.customersMutex.RLock()
    defer s.customersMutex.RUnlock()

    // Convert map of customers into array
    values := make([]Customer, len(s.customers))
    i := 0
    for _, v := range s.customers {
        values[i] = v
        i++
    }

    // Sort result if orderBy is present
    if len(orderBy) > 0 {
        if orderBy != "companyName" {
            return nil, ErrInvalidOrderBy
        }

        sort.Sort(ByCompanyName(values))
    }

    return values, nil
}

func (s customerRepository) GetCustomer(ctx context.Context, cid uuid.UUID) (Customer, error) {
    // Lock customers while accessing it
    s.customersMutex.RLock()
    defer s.customersMutex.RUnlock()

    // Check if customer with given ID exists
    if c, ok := s.customers[cid]; ok {
        return c, nil
    }

    return Customer{}, ErrNotFound
}

func (s customerRepository) AddCustomer(ctx context.Context, c Customer) (Customer, error) {
    // Lock customers while accessing it
    s.customersMutex.Lock()
    defer s.customersMutex.Unlock()

    // Make sure that incoming custer data is sane
    if c.CustomerID != uuid.Nil {
        return Customer{}, ErrGivenCustomerID
    }

    if len(c.CompanyName) == 0 {
        return Customer{}, ErrMissingMandatoryValue{Field: "CompanyName"}
    }

    if len(c.ContactName) == 0 {
        return Customer{}, ErrMissingMandatoryValue{Field: "ContactName"}
    }

    if len(c.Country) != 3 {
        return Customer{}, ErrInvalidCountry
    }

    if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        return Customer{}, ErrInvalidHourlyRate
    }

    // Assign new customer ID
    c.CustomerID, _ = uuid.NewUUID()

    // Add customer to our list
    s.customers[c.CustomerID] = c

    return c, nil
}

func (s customerRepository) DeleteCustomer(ctx context.Context, cid uuid.UUID) error {
    // Lock customers while accessing it
    s.customersMutex.Lock()
    defer s.customersMutex.Unlock()

    // Check if customer with given ID exists
    if _, ok := s.customers[cid]; ok {
        delete(s.customers, cid)
        return nil
    }

    return ErrNotFound
}

func (s customerRepository) PatchCustomer(ctx context.Context, cid uuid.UUID, c Customer) (Customer, error) {
    // Lock customers while accessing it
    s.customersMutex.Lock()
    defer s.customersMutex.Unlock()

    // Check if customer with given ID exists
    if cOld, ok := s.customers[cid]; ok {
        // Update specified fields
        if len(c.CompanyName) > 0 {
            cOld.CompanyName = c.CompanyName
        }

        if len(c.ContactName) > 0 {
            cOld.ContactName = c.ContactName
        }

        if len(c.Country) > 0 {
            cOld.Country = c.Country
        }

        if c.HourlyRate != decimal.NewFromInt(0) {
            cOld.HourlyRate = c.HourlyRate
        }

        // Update customer in in-memory store
        s.customers[cid] = cOld

        return cOld, nil
    }

    return Customer{}, ErrNotFound
}

// ByCompanyName is used for sorting customers by company name
type ByCompanyName []Customer

func (c ByCompanyName) Len() int           { return len(c) }
func (c ByCompanyName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCompanyName) Less(i, j int) bool { return c[i].CompanyName < c[j].CompanyName }
```

## Add Endpoints

* Create *customerendpoints.go*

```go
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
```

## Add Transport

* Create *customertransport.go*

```go
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
```

## Add Logging Middleware

* Create *customermiddlewares.go*

```go
package customersvc

import (
    "context"
    "time"

    "github.com/go-kit/kit/log"
    "github.com/google/uuid"
)

// CustomerMiddleware wraps a CustomerService to surround business logic methods
// with additional functionalities like e.g. logging
type CustomerMiddleware func(CustomerService) CustomerService

// CustomerLoggingMiddleware returns a factory for a logging middleware for the customer service.
func CustomerLoggingMiddleware(logger log.Logger) CustomerMiddleware {
    return func(next CustomerService) CustomerService {
        return &customerLoggingMiddleware{
            next:   next,
            logger: logger,
        }
    }
}

type customerLoggingMiddleware struct {
    next   CustomerService
    logger log.Logger
}

func (mw customerLoggingMiddleware) GetCustomers(ctx context.Context, orderBy string) (c []Customer, err error) {
    defer func(begin time.Time) {
        mw.logger.Log("method", "GetCustomers", "took", time.Since(begin), "err", err)
    }(time.Now())
    return mw.next.GetCustomers(ctx, orderBy)
}

func (mw customerLoggingMiddleware) GetCustomer(ctx context.Context, cid uuid.UUID) (c Customer, err error) {
    defer func(begin time.Time) {
        mw.logger.Log("method", "GetCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
    }(time.Now())
    return mw.next.GetCustomer(ctx, cid)
}

func (mw customerLoggingMiddleware) AddCustomer(ctx context.Context, c Customer) (cust Customer, err error) {
    defer func(begin time.Time) {
        mw.logger.Log("method", "AddCustomer", "id", c.CustomerID.String(), "took", time.Since(begin), "err", err)
    }(time.Now())
    return mw.next.AddCustomer(ctx, c)
}

func (mw customerLoggingMiddleware) DeleteCustomer(ctx context.Context, cid uuid.UUID) (err error) {
    defer func(begin time.Time) {
        mw.logger.Log("method", "DeleteCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
    }(time.Now())
    return mw.next.DeleteCustomer(ctx, cid)
}

func (mw customerLoggingMiddleware) PatchCustomer(ctx context.Context, cid uuid.UUID, c Customer) (cust Customer, err error) {
    defer func(begin time.Time) {
        mw.logger.Log("method", "PatchCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
    }(time.Now())
    return mw.next.PatchCustomer(ctx, cid, c)
}
```

## Add Main Method

* Create folder *cmd/server*
  
* Create *main.go* in the previously created folder.

```go
package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/go-kit/kit/log"
    customersvc "github.com/rstropek/golang-samples/go-kit"
)

func main() {
    var (
        httpAddr = flag.String("http.addr", ":4000", "HTTP listen address")
    )
    flag.Parse()

    // Create logger
    var logger log.Logger
    {
        logger = log.NewLogfmtLogger(os.Stderr)
        logger = log.With(logger, "ts", log.DefaultTimestampUTC)
        logger = log.With(logger, "caller", log.DefaultCaller)
    }

    // Create customer service and surround it with logging middleware
    var s customersvc.CustomerService
    {
        s = customersvc.NewCustomerRepository()
        s = customersvc.CustomerLoggingMiddleware(logger)(s)
    }

    // Create HTTP transport for customer service
    var h http.Handler
    {
        h = customersvc.MakeCustomerHTTPHandler(s, log.With(logger, "component", "HTTP"))
    }

    errs := make(chan error)
    go func() {
        c := make(chan os.Signal)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errs <- fmt.Errorf("%s", <-c)
    }()

    go func() {
        logger.Log("transport", "HTTP", "addr", *httpAddr)
        errs <- http.ListenAndServe(*httpAddr, h)
    }()

    logger.Log("exit", <-errs)
}
```

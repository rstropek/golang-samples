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

	// Make sure that incoming customer data is sane
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

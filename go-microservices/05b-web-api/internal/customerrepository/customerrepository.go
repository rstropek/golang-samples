package customerrepository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Customer holds data of a customer record
type Customer struct {
	CustomerID  uuid.UUID       `json:"customerID,omitempty"`
	CompanyName string          `json:"customerName"`
	ContactName string          `json:"contactName"`
	Country     string          `json:"country"`
	HourlyRate  decimal.Decimal `json:"hourlyRate"`
}

// CustomerRepository is an in-memory repository of customers
type CustomerRepository struct {
	// Store map of customers in memory
	customers map[uuid.UUID]Customer

	// Mutex serializing access to customers. We need this mutex because
	// go serves all incoming HTTP requests in their own goroutine. Therefore,
	// it is possible if not likely that handlers will run concurrently.
	// As concurrent reading without writing is allowed, we could optimize
	// our code using `RWMutex` (https://golang.org/pkg/sync/#RWMutex).
	// However, this is out of scope for this sample.
	customersMutex *sync.Mutex
}

// NewCustomerRepository creates a customer repository
func NewCustomerRepository() CustomerRepository {
	return CustomerRepository{
		customers:      make(map[uuid.UUID]Customer, 0),
		customersMutex: &sync.Mutex{},
	}
}

// GetCustomerByID looks for a customer with a given ID
func (cr CustomerRepository) GetCustomerByID(cid uuid.UUID) (*Customer, bool) {
	// Lock customers while accessing it
	cr.customersMutex.Lock()
	defer cr.customersMutex.Unlock()

	// Check if customer with given ID exists
	if c, ok := cr.customers[cid]; ok {
		return &c, true
	}

	return nil, false
}

// GetCustomersArray returns all stored customers as an array
func (cr CustomerRepository) GetCustomersArray() []Customer {
	// Lock customers while accessing it
	cr.customersMutex.Lock()
	defer cr.customersMutex.Unlock()

	// Convert map of customers into array
	values := make([]Customer, len(cr.customers))
	i := 0
	for _, v := range cr.customers {
		values[i] = v
		i++
	}

	return values
}

// AddCustomer adds a customer to the repository
func (cr CustomerRepository) AddCustomer(c Customer) {
	// Lock customers while accessing it
	cr.customersMutex.Lock()
	defer cr.customersMutex.Unlock()

	// Add customer to our list
	cr.customers[c.CustomerID] = c
}

// DeleteCustomerByID removes a customer with a given ID
func (cr CustomerRepository) DeleteCustomerByID(cid uuid.UUID) bool {
	// Lock customers while accessing it
	cr.customersMutex.Lock()
	defer cr.customersMutex.Unlock()

	// Check if customer with given ID exists
	if _, ok := cr.customers[cid]; ok {
		delete(cr.customers, cid)
		return true
	}

	return false
}

// PatchCustomer patches a customer with the given values
func (cr CustomerRepository) PatchCustomer(cid uuid.UUID, c Customer) (*Customer, bool) {
	// Lock customers while accessing it
	cr.customersMutex.Lock()
	defer cr.customersMutex.Unlock()

	// Check if customer with given ID exists
	if cOld, ok := cr.customers[cid]; ok {
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
		cr.customers[cid] = cOld

		return &cOld, true
	}

	return nil, false
}

// ByCompanyName is used for sorting customers by company name
type ByCompanyName []Customer

func (c ByCompanyName) Len() int           { return len(c) }
func (c ByCompanyName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCompanyName) Less(i, j int) bool { return c[i].CompanyName < c[j].CompanyName }

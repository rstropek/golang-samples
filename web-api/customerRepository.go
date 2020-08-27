package main

import (
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Setup structure for storing customer data
type customer struct {
	CustomerID  uuid.UUID       `json:"customerID,omitempty"`
	CompanyName string          `json:"customerName"`
	ContactName string          `json:"contactName"`
	Country     string          `json:"country"`
	HourlyRate  decimal.Decimal `json:"hourlyRate"`
}

// Store map of customers in memory
var customers = make(map[uuid.UUID]customer, 0)

// Mutex serializing access to customers
var customersMutex = &sync.Mutex{}

// getCustomersArray returns all stored customers as an array
func getCustomersArray() []customer {
	// Lock customers while accessing it
	customersMutex.Lock()
	defer customersMutex.Unlock()

	// Convert map of customers into array
	values := make([]customer, len(customers))
	i := 0
	for _, v := range customers {
		values[i] = v
		i++
	}

	return values
}

type byCompanyName []customer

func (c byCompanyName) Len() int           { return len(c) }
func (c byCompanyName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCompanyName) Less(i, j int) bool { return c[i].CompanyName < c[j].CompanyName }

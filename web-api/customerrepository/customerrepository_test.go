package customerrepository

import (
	"sort"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestAddCustomer(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{})

	assert.Equal(t, 1, len(cr.customers))
}

func TestGetCustomersArray(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{})

	assert.Equal(t, 1, len(cr.GetCustomersArray()))
}

func TestGetCustomerByID(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{CustomerID: uuid.Nil})

	_, ok := cr.GetCustomerByID(uuid.Nil)
	assert.True(t, ok)
}

func TestDeleteCustomerByID(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{CustomerID: uuid.Nil})

	cr.DeleteCustomerByID(uuid.Nil)
	assert.Equal(t, 0, len(cr.customers))
}

func TestPatchCustomer(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{
		CustomerID:  uuid.Nil,
		CompanyName: "Acme Corp",
	})

	cr.PatchCustomer(uuid.Nil, Customer{CompanyName: "Foo Bar"})
	assert.Equal(t, "Foo Bar", cr.customers[uuid.Nil].CompanyName)
}

func TestOrderByCompanyName(t *testing.T) {
	cr := NewCustomerRepository()
	cr.AddCustomer(Customer{CompanyName: "B"})
	cr.AddCustomer(Customer{CompanyName: "A"})

	c := cr.GetCustomersArray()
	sort.Sort(ByCompanyName(c))
	assert.Equal(t, "A", c[0].CompanyName)
}

# Storyboard

## Getting started

* Create empty directory *basicwebapi*

```bash
touch main.go
go mod init github.com/rstropek/golang-samples/basicwebapi
```

* Add starter code to *main.go*

```go
package main

import (
    "log"
    "net/http"
)

// Define a home handler function which writes a byte slice containing
// hard-coded JSON as the response body.
func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{ \"foo\": \"bar\" }"))
}

func main() {
    // Initialize a new servemux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)

    // Use the http.ListenAndServe() function to start a new web server.
    port := ":4000"
    log.Printf("Starting server on %s", port)
    err := http.ListenAndServe(port, mux)
    log.Fatal(err)
}
```

* Run app

```bash
go run .
```

* Test it

```http
GET http://localhost:4000/
```

## Add Customer

* Add package for handling GUIDs and decimal values

```bash
go get github.com/google/uuid
go get github.com/shopspring/decimal
```

* Add customer struct

```go
// ...

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// ...

// Setup structure for storing customer data
type customer struct {
    CustomerID  uuid.UUID       `json:"customerID,omitempty"`
    CompanyName string          `json:"customerName"`
    ContactName string          `json:"contactName"`
    Country     string          `json:"country"`
    HourlyRate  decimal.Decimal `json:"hourlyRate"`
}
```

* Change home function to return object encoded in JSON

```go
// Return encoded demo customer in JSON
func home(w http.ResponseWriter, r *http.Request) {
    cid, _ := uuid.NewUUID()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&customer {
        CustomerID: cid,
        CompanyName: "Acme Corp",
        ContactName: "Foo Bar",
        Country: "DEU",
        HourlyRate: decimal.NewFromInt(42),
    })
}
```

* Test it

```http
GET http://localhost:4000/
```

## Add More Powerful Router

* Add *Gorilla MUX* package

```bash
go get github.com/gorilla/mux
```

* Change mux to Gorilla

```go
// ...

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

// ...

func main() {
    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/", home)

    // ...
}
```

* Test it

```http
GET http://localhost:4000/
```

## Store Customers in In-Memory Map

* Remove `home` method

```go
// ...
import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

// ...

// Store map of customers in memory
var customers = make(map[uuid.UUID]customer, 0)

// Mutex serializing access to customers. We need this mutex because
// go serves all incoming HTTP requests in their own goroutine. Therefore,
// it is possible if not likely that handlers will run concurrently.
// As concurrent reading without writing is allowed, we could optimize
// our code using `RWMutex` (https://golang.org/pkg/sync/#RWMutex).
// However, this is out of scope for this sample.
var customersMutex = &sync.Mutex{}

// ...

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

func getCustomers(w http.ResponseWriter, r *http.Request) {
    // Return all customers
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(getCustomersArray())
}

// ...

func main() {
    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")

    // Use the http.ListenAndServe() function to start a new web server.
    port := ":4000"
    log.Printf("Starting server on %s", port)
    err := http.ListenAndServe(port, mux)
    log.Fatal(err)
}
```

## Command-Line Arguments

```go
// ...

import (
    "encoding/json"
    "flag"
    "log"
    "net/http"
    "sync"
    "fmt"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

// ...

func main() {
    // Parse command-line arguments
    var portFlag = flag.Uint("p", 4000, "Port number for starting server")
    flag.Parse()

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")

    // Use the http.ListenAndServe() function to start a new web server.
    log.Printf("Starting server on %d", *portFlag)
    err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), mux)
    log.Fatal(err)
}
```

* Test it: `go run . -p 8081`

## Get Single Customer

```go
// ...

func getCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        // Note http.Error shortcut. Use it to send a non-200 status code and
        // plain-text response body.
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if c, ok := customers[cid]; ok {
        // Return customer
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(c)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

// newUUID returns a new UUID and ignores potential errors
func newUUID() uuid.UUID {
    r, _ := uuid.NewUUID()
    return r
}

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")

    // ...
}
```

* Test it

```http
# @name customers
GET http://localhost:4000/customers

###

@customerID = {{customers.response.body.$[0].customerID}}

GET http://localhost:4000/customers/{{customerID}}

###
GET http://localhost:4000/customers/00000000-0000-0000-0000-000000000000
```

## Add Customer

```go
// ...

func addCustomer(w http.ResponseWriter, r *http.Request) {
    // Decode customer data from request body
    var c = customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // Make sure that incoming custer data is sane
    if c.CustomerID != uuid.Nil {
        http.Error(w, "CustomerID must be empty", http.StatusBadRequest)
        return
    }

    if len(c.CompanyName) == 0 {
        http.Error(w, "Company name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.ContactName) == 0 {
        http.Error(w, "Contact name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    // Assign new customer ID
    c.CustomerID = newUUID()

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Add customer to our list
    customers[c.CustomerID] = c

    // Return customer
    w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", addCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")

    // ...
}
```

* Test it

```http
# ...

###
POST http://localhost:4000/customers

{
    "customerName": "Acme Corp",
    "contactName": "Foo Bar",
    "country": "DEU",
    "hourlyRate": "42"
}
```

## Delete Customer

```go
// ...


func deleteCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if _, ok := customers[cid]; ok {
        delete(customers, cid)
        w.WriteHeader(http.StatusNoContent)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", addCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", deleteCustomer).Methods("DELETE")

    // ...
}
```

```http
# ...

###
DELETE http://localhost:4000/customers/{{customerID}}
```

## Update Customer

```go
// ...

func patchCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Decode customer data from request body
    var c = customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // If customer ID was specified, it must match the customer ID from path
    if c.CustomerID != uuid.Nil && cid != c.CustomerID {
        http.Error(w, "Cannot update customer ID", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if cOld, ok := customers[cid]; ok {
        // Update specified fields
        if len(c.CompanyName) > 0 {
            cOld.CompanyName = c.CompanyName
        }

        if len(c.ContactName) > 0 {
            cOld.ContactName = c.ContactName
        }

        if len(c.Country) > 0 {
            if len(c.Country) != 3 {
                http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
                return
            }

            cOld.Country = c.Country
        }

        if c.HourlyRate != decimal.NewFromInt(0) {
            if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
                http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
                return
            }

            cOld.HourlyRate = c.HourlyRate
        }

        // Update customer in in-memory store
        customers[cid] = cOld

        // Return updated customer data
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cOld)
    }
}

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", addCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", deleteCustomer).Methods("DELETE")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", patchCustomer).Methods("PATCH")

    // ...
}
```

## Add Query Parameter

```go
// ...

import (
    // ...
    "sort"
    // ...
)

// ...

type byCompanyName []customer

func (c byCompanyName) Len() int           { return len(c) }
func (c byCompanyName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCompanyName) Less(i, j int) bool { return c[i].CompanyName < c[j].CompanyName }

func getCustomers(w http.ResponseWriter, r *http.Request) {
    custArray := getCustomersArray()
    orderBy := r.FormValue("orderBy")
    if len(orderBy) > 0 {
        if orderBy != "companyName" {
            http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
            return
        }

        sort.Sort(byCompanyName(custArray))
    }

    // Return all customers
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(custArray)
}

// ...

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) { panic("Something really bad happened...") }).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Queries("orderBy", "{orderBy}").Methods("GET")

    // ...
}
```

* Try it

```http
# ...

###
GET http://localhost:4000/customers?orderBy=companyName
```

## Add Middleware

* Add *negroni*

```bash
go get github.com/urfave/negroni
go get github.com/rs/cors
```

* Add classic middleware and CORS

```go
// ...
import (
    // ...

    "github.com/urfave/negroni"
    "github.com/rs/cors"
)

// ...

func main() {
    // ...

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) { panic("Something really bad happened...") }).Methods("GET")
    // ...

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Use(cors.AllowAll())

    // Use the http.ListenAndServe() function to start a new web server.
    log.Printf("Starting server on %d", *portFlag)
    err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), n)
    log.Fatal(err)
}
```

* Try it

```http
# ...

###
GET http://localhost:4000/panic
```

* Create *public* subdirectory

* Add demo client

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Demo Client</title>
</head>
<body>
    <ul id="customers" />

    <script>
        (async () => {
            const cElem = document.getElementById("customers");

            let html = "";
            const result = await fetch("/customers");
            const custs = await result.json();
            for (const c of custs) {
                html += `<li>${c.customerName}</li>`;
            }

            cElem.innerHTML = html;
        })();
    </script>
</body>
</html>
```

* Try it by opening `http://localhost:4000/index.html` in your browser

## Split Into Multiple Files

* Create *customerrepository.go*

```go
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
```

* Create *handlers.go*

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/shopspring/decimal"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
    custArray := getCustomersArray()
    orderBy := r.FormValue("orderBy")
    if len(orderBy) > 0 {
        if orderBy != "companyName" {
            http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
            return
        }

        sort.Sort(byCompanyName(custArray))
    }

    // Return all customers
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(custArray)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if c, ok := customers[cid]; ok {
        // Return customer
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(c)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

// newUUID returns a new UUID and ignores potential errors
func newUUID() uuid.UUID {
    r, _ := uuid.NewUUID()
    return r
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
    // Decode customer data from request body
    var c = customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // Make sure that incoming custer data is sane
    if c.CustomerID != uuid.Nil {
        http.Error(w, "CustomerID must be empty", http.StatusBadRequest)
        return
    }

    if len(c.CompanyName) == 0 {
        http.Error(w, "Company name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.ContactName) == 0 {
        http.Error(w, "Contact name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    // Assign new customer ID
    c.CustomerID = newUUID()

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Add customer to our list
    customers[c.CustomerID] = c

    // Return customer
    w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if _, ok := customers[cid]; ok {
        delete(customers, cid)
        w.WriteHeader(http.StatusNoContent)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

func patchCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Decode customer data from request body
    var c = customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // If customer ID was specified, it must match the customer ID from path
    if c.CustomerID != uuid.Nil && cid != c.CustomerID {
        http.Error(w, "Cannot update customer ID", http.StatusBadRequest)
        return
    }

    // Lock customers while accessing it
    customersMutex.Lock()
    defer customersMutex.Unlock()

    // Check if customer with given ID exists
    if cOld, ok := customers[cid]; ok {
        // Update specified fields
        if len(c.CompanyName) > 0 {
            cOld.CompanyName = c.CompanyName
        }

        if len(c.ContactName) > 0 {
            cOld.ContactName = c.ContactName
        }

        if len(c.Country) > 0 {
            if len(c.Country) != 3 {
                http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
                return
            }

            cOld.Country = c.Country
        }

        if c.HourlyRate != decimal.NewFromInt(0) {
            if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
                http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
                return
            }

            cOld.HourlyRate = c.HourlyRate
        }

        // Update customer in in-memory store
        customers[cid] = cOld

        // Return updated customer data
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cOld)
    }
}
```

* Shorten *main.go*

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "github.com/shopspring/decimal"
    "github.com/urfave/negroni"
)

func main() {
    // Parse command-line arguments
    var portFlag = flag.Uint("p", 4000, "Port number for starting server")
    flag.Parse()

    // Add one demo record
    cid := newUUID()
    customers[cid] = customer{
        CustomerID:  cid,
        CompanyName: "Acme Corp",
        ContactName: "Foo Bar",
        Country:     "DEU",
        HourlyRate:  decimal.NewFromInt(42),
    }

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) { panic("Something really bad happened...") }).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Queries("orderBy", "{orderBy}").Methods("GET")
    mux.HandleFunc("/customers", addCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", deleteCustomer).Methods("DELETE")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", patchCustomer).Methods("PATCH")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Use(cors.AllowAll())

    // Use the http.ListenAndServe() function to start a new web server.
    log.Printf("Starting server on %d", *portFlag)
    err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), n)
    log.Fatal(err)
}
```

## Convert Customer Repository in Package

* Create subfolder *customerrepository*

* Move *customerrepository.go* into new folder

* Change code of *customerrepository.go*

```go
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
```

* Adjust code in *handlers.go*

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/rstropek/golang-samples/web-api/customerrepository"
    "github.com/shopspring/decimal"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
    custArray := repo.GetCustomersArray()
    orderBy := r.FormValue("orderBy")
    if len(orderBy) > 0 {
        if orderBy != "companyName" {
            http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
            return
        }

        sort.Sort(customerrepository.ByCompanyName(custArray))
    }

    // Return all customers
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(custArray)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Check if customer with given ID exists
    if c, ok := repo.GetCustomerByID(cid); ok {
        // Return customer
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(c)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

// newUUID returns a new UUID and ignores potential errors
func newUUID() uuid.UUID {
    r, _ := uuid.NewUUID()
    return r
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
    // Decode customer data from request body
    var c = customerrepository.Customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // Make sure that incoming custer data is sane
    if c.CustomerID != uuid.Nil {
        http.Error(w, "CustomerID must be empty", http.StatusBadRequest)
        return
    }

    if len(c.CompanyName) == 0 {
        http.Error(w, "Company name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.ContactName) == 0 {
        http.Error(w, "Contact name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    // Assign new customer ID
    c.CustomerID = newUUID()

    // Add customer to our list
    repo.AddCustomer(c)

    // Return customer
    w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Delete customer
    if repo.DeleteCustomerByID(cid) {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

func patchCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Decode customer data from request body
    var c = customerrepository.Customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // If customer ID was specified, it must match the customer ID from path
    if c.CustomerID != uuid.Nil && cid != c.CustomerID {
        http.Error(w, "Cannot update customer ID", http.StatusBadRequest)
        return
    }

    if len(c.Country) > 0 && len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if c.HourlyRate != decimal.NewFromInt(0) && decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    if cNew, ok := repo.PatchCustomer(cid, c); ok {
        // Return updated customer data
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cNew)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}
```

* Adjust code in *main.go*

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "github.com/rstropek/golang-samples/web-api/customerrepository"
    "github.com/shopspring/decimal"
    "github.com/urfave/negroni"
)

var repo = customerrepository.NewCustomerRepository()

func main() {
    // Parse command-line arguments
    var portFlag = flag.Uint("p", 4000, "Port number for starting server")
    flag.Parse()

    // Add one demo record
    cid := newUUID()
    repo.AddCustomer(customerrepository.Customer{
        CustomerID:  cid,
        CompanyName: "Acme Corp",
        ContactName: "Foo Bar",
        Country:     "DEU",
        HourlyRate:  decimal.NewFromInt(42),
    })

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) { panic("Something really bad happened...") }).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Methods("GET")
    mux.HandleFunc("/customers", getCustomers).Queries("orderBy", "{orderBy}").Methods("GET")
    mux.HandleFunc("/customers", addCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", getCustomer).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", deleteCustomer).Methods("DELETE")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", patchCustomer).Methods("PATCH")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Use(cors.AllowAll())

    // Use the http.ListenAndServe() function to start a new web server.
    log.Printf("Starting server on %d", *portFlag)
    err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), n)
    log.Fatal(err)
}
```

* Try it

## Add Unit Tests For Customer Repository

* Add *customerrepository_test.go*

```go
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
```

* Try it: `go test .`

## Convert Customer API Handlers in Package

* Create subfolder *customerhandlers*

* Move *customerhandlers.go* into new folder

* Change code of *customerhandlers.go*

```go
package customerhandlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/rstropek/golang-samples/web-api/customerrepository"
    "github.com/shopspring/decimal"
)

// CustomerHandlers represents functions handling HTTP requests for customers management web api
type CustomerHandlers struct {
    repo customerrepository.CustomerRepository
}

// NewCustomerHandlers creates a customer handler object
func NewCustomerHandlers(repo customerrepository.CustomerRepository) CustomerHandlers {
    return CustomerHandlers{
        repo: repo,
    }
}

// GetCustomers returns all customers
func (ch CustomerHandlers) GetCustomers(w http.ResponseWriter, r *http.Request) {
    custArray := ch.repo.GetCustomersArray()
    orderBy := r.FormValue("orderBy")
    if len(orderBy) > 0 {
        if orderBy != "companyName" {
            http.Error(w, "Currently, we can only order by companyName", http.StatusBadRequest)
            return
        }

        sort.Sort(customerrepository.ByCompanyName(custArray))
    }

    // Return all customers
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(custArray)
}

// GetCustomer returns a single customer based on a given customer ID
func (ch CustomerHandlers) GetCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Check if customer with given ID exists
    if c, ok := ch.repo.GetCustomerByID(cid); ok {
        // Return customer
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(c)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

// newUUID returns a new UUID and ignores potential errors
func newUUID() uuid.UUID {
    r, _ := uuid.NewUUID()
    return r
}

// AddCustomer adds a customer
func (ch CustomerHandlers) AddCustomer(w http.ResponseWriter, r *http.Request) {
    // Decode customer data from request body
    var c = customerrepository.Customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // Make sure that incoming custer data is sane
    if c.CustomerID != uuid.Nil {
        http.Error(w, "CustomerID must be empty", http.StatusBadRequest)
        return
    }

    if len(c.CompanyName) == 0 {
        http.Error(w, "Company name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.ContactName) == 0 {
        http.Error(w, "Contact name must not be empty", http.StatusBadRequest)
        return
    }

    if len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    // Assign new customer ID
    c.CustomerID = newUUID()

    // Add customer to our list
    ch.repo.AddCustomer(c)

    // Return customer
    w.Header().Set("Location", fmt.Sprintf("/customers/%s", c.CustomerID))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

// DeleteCustomer deletes a customer based on a given ID
func (ch CustomerHandlers) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Delete customer
    if ch.repo.DeleteCustomerByID(cid) {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}

// PatchCustomer patches a customer based on a given ID and new field values
func (ch CustomerHandlers) PatchCustomer(w http.ResponseWriter, r *http.Request) {
    // Get customer ID from path
    cid, err := uuid.Parse(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
        return
    }

    // Decode customer data from request body
    var c = customerrepository.Customer{}
    if json.NewDecoder(r.Body).Decode(&c) != nil {
        http.Error(w, "Could not deserialize customer from HTTP body", http.StatusBadRequest)
        return
    }

    // If customer ID was specified, it must match the customer ID from path
    if c.CustomerID != uuid.Nil && cid != c.CustomerID {
        http.Error(w, "Cannot update customer ID", http.StatusBadRequest)
        return
    }

    if len(c.Country) > 0 && len(c.Country) != 3 {
        http.Error(w, "Country name must be three characters long (use ISO 3166-1 Alpha-3 code)", http.StatusBadRequest)
        return
    }

    if c.HourlyRate != decimal.NewFromInt(0) && decimal.NewFromInt(0).GreaterThan(c.HourlyRate) {
        http.Error(w, "Hourly rate must be >= 0", http.StatusBadRequest)
        return
    }

    if cNew, ok := ch.repo.PatchCustomer(cid, c); ok {
        // Return updated customer data
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cNew)
        return
    }

    // Customer hasn't been found
    http.NotFound(w, r)
}
```

* Adjust code in *main.go*

```go
package main

import (
    "github.com/google/uuid"
    "github.com/rstropek/golang-samples/web-api/customerhandlers"
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "github.com/rstropek/golang-samples/web-api/customerrepository"
    "github.com/shopspring/decimal"
    "github.com/urfave/negroni"
)


func main() {
    // Parse command-line arguments
    var portFlag = flag.Uint("p", 4000, "Port number for starting server")
    flag.Parse()

    // Create customer repository
    repo := customerrepository.NewCustomerRepository()

    // Add one demo record
    cid, _ := uuid.NewUUID()
    repo.AddCustomer(customerrepository.Customer{
        CustomerID:  cid,
        CompanyName: "Acme Corp",
        ContactName: "Foo Bar",
        Country:     "DEU",
        HourlyRate:  decimal.NewFromInt(42),
    })

    // Create handlers
    ch := customerhandlers.NewCustomerHandlers(repo)

    // Initialize a new Gorilla mux, then register the home function as
    // the handler for the "/" URL pattern.
    mux := mux.NewRouter()
    mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) { panic("Something really bad happened...") }).Methods("GET")
    mux.HandleFunc("/customers", ch.GetCustomers).Methods("GET")
    mux.HandleFunc("/customers", ch.GetCustomers).Queries("orderBy", "{orderBy}").Methods("GET")
    mux.HandleFunc("/customers", ch.AddCustomer).Methods("POST")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", ch.GetCustomer).Methods("GET")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", ch.DeleteCustomer).Methods("DELETE")
    mux.HandleFunc("/customers/{id:[0-9A-Fa-f]{8}(?:-[0-9A-Fa-f]{4}){3}-[0-9A-Fa-f]{12}}", ch.PatchCustomer).Methods("PATCH")

    n := negroni.Classic()
    n.UseHandler(mux)
    n.Use(cors.AllowAll())

    // Use the http.ListenAndServe() function to start a new web server.
    log.Printf("Starting server on %d", *portFlag)
    err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), n)
    log.Fatal(err)
}
```

## Add Sample Unit Test for Customer API Handlers

* Add *customerhandlers_test.go*

```go
package customerhandlers

import (
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "github.com/rstropek/golang-samples/web-api/customerrepository"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGetCustomers(t *testing.T) {
    // Here we use the existing customer repository. In practice, you would probably
    // use a mocking framework like https://github.com/stretchr/testify. However, proper
    // mocking for unit tests is out of scope here.
    repo := customerrepository.NewCustomerRepository()
    repo.AddCustomer(customerrepository.Customer{CompanyName: "Foo Bar"})
    ch := NewCustomerHandlers(repo)

    // Create a request to pass to our handler
    req, _ := http.NewRequest("GET", "/", nil)

    // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ch.GetCustomers)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    assert.Equal(t, http.StatusOK, rr.Code)

    // Check content type
    assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

    // Check the JSON result
    result := make([]customerrepository.Customer, 0)
    json.NewDecoder(rr.Body).Decode(&result)
    assert.Equal(t, 1, len(result))
    assert.Equal(t, "Foo Bar", result[0].CompanyName)
}
```

* Try it: `go test .`

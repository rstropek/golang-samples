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

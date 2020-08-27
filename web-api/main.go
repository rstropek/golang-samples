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

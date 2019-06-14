package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
)

func ping(w http.ResponseWriter, r *http.Request) {
	// Note type conversion from string to byte array
	w.Write([]byte("pong"))
}

func main() {
	// Note the use of https://golang.org/pkg/flag/ for command-line parsing
	port := flag.Uint("p", 8080, "Port on which the server should listen")
	flag.Parse()

	// For demo purposes, we create a custom request MUX (multiplexer).
	// To keep things simple, we could have used the DefaultServeMux
	// (see also https://golang.org/pkg/net/http/#Handle)
	mux := http.NewServeMux()

	// Handle request with a function
	mux.HandleFunc("/ping", ping)

	// Handle request with an anonymous function
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		// Note how we access query parameters
		msg := r.URL.Query().Get("msg")
		if len(msg) > 0 {
			w.Write([]byte(msg))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// Use an existing handler for static files
	fs := http.FileServer(http.Dir(path.Join("client", ".")))
	mux.Handle("/client/", http.StripPrefix("/client", fs))

	// Start web server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	// Create a channel (see also https://gobyexample.com/channels) that
	// receives a signal from the OS if the process should be ended
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Setup HTTP server
	http.HandleFunc("/ping", ping)
	h := &http.Server{Addr: ":8080"}

	// Start web server on a separate goroutine
	go func() {
		fmt.Println("Listening on http://0.0.0.0:8080")
		log.Fatal(h.ListenAndServe())
	}()

	// Wait until we receive a sto signal
	<-stop

	// Shutdown server (note that this does not interrupt active connections)
	fmt.Println("\nShutting down the server...")
	h.Shutdown(context.Background())
	fmt.Println("Server gracefully stopped")
}

package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var errorLogger *log.Logger
var infoLogger *log.Logger

func main() {
	// Create loggers for error and info messages
	errorLogger = log.New(os.Stderr, "Echo bot", log.LstdFlags)
	infoLogger = log.New(os.Stdout, "Echo bot", log.LstdFlags)

	// Create a listener for incoming connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		// Listener creation failed, log the error and exit
		errorLogger.Println("Error listening:", err)
		return
	}
	// Close the listener when the application closes
	defer listener.Close()

	infoLogger.Println("Server started, waiting for incoming connections...")

	// Create a context that can be used to cancel the server
	ctx, cancel := context.WithCancel(context.Background())

	// Create a wait group to wait for all connections to be closed
	var wg sync.WaitGroup

	// Create a channel that will receive a signal when the server
	// should be stopped (e.g. when the user presses Ctrl+C)
	quit := CreateSigintQuitChannel()

	// We need the listener's underlying TCP listener to set a deadline
	tcpListener, _ := listener.(*net.TCPListener)

Accept:
	for {
		// Set a deadline for the Accept() call so that we can check
		// if the server is shutting down
		tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
		conn, err := listener.Accept()
		if err != nil {
			if IsTimeout(err) {
				select {
				case <-quit:
					// We received a signal to stop the server.
					// Cancel the context and wait for all connections
					// to be closed.
					infoLogger.Println("Server is shutting down...")
					cancel()
					wg.Wait()

					// Stop accepting new connections
					break Accept
				default:
					// No quit signal received, continue accepting connections
					continue Accept
				}
			}
			
			// Error accepting connection, log the error and stop 
			// accepting new connections
			errorLogger.Println("Error accepting connection:", err)
			break Accept
		}
		
		// No error happend, handle the connection in a separate goroutine
		wg.Add(1)
		go handleConnection(ctx, &wg, conn)
	}

	infoLogger.Println("Server stopped")
}

func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	// Close the connection when the function returns
	defer wg.Done()
	defer conn.Close()

	infoLogger.Println("Received connection from", conn.RemoteAddr().String())

	// Create a buffered reader to read messages from the connection
	reader := bufio.NewReader(conn)
Messages:
	for {
		select {
		case <-ctx.Done():
			// The server is shutting down, stop handling messages
			break Messages
		default:
			// Set a deadline for the ReadString() call so that we can
			// regularly check if the server is shutting down
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			message, err := reader.ReadString('\n')
			if err != nil {
				if IsTimeout(err) {
					// No message received, continue reading
					continue Messages
				}

				// Error reading from connection, log the error and stop
				// handling messages
				errorLogger.Println("Error reading from connection:", err)
				break Messages
			}

			// Print message to the console
			message = strings.TrimSpace(message)
			if message == "" {
				continue
			}
			infoLogger.Println("Received message:", message)

			// Write the message back to the connection
			_, err = conn.Write([]byte(message + "\n"))
			if err != nil {
				// Error writing to connection, log the error and stop
				// handling messages
				errorLogger.Println("Error writing to connection:", err)
				break Messages
			}
		}
	}

	infoLogger.Println("Connection closed")
}

func IsTimeout(err error) bool {
	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		return true
	}
	return false
}

func CreateSigintQuitChannel() <-chan any {
	// Create a channel to listen for SIGINT and SIGTERM signals.
	// When a signal is received, we signal the quit channel to
	// stop the server.
	sig := make(chan os.Signal, 1)
	quit := make(chan any, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		quit <- sig
	}()
	return quit
}
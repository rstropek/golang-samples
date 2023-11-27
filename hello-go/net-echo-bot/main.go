package main

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var logger *slog.Logger

func main() {
	// Create loggers for error and info messages
	// logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create a listener for incoming connections
	const PORT = ":8080"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		// Listener creation failed, log the error and exit
		logger.Error("error listening", "error details", err)
		return
	}
	// Close the listener when the application closes
	defer listener.Close()

	logger.Info("server started, waiting for incoming connections", "port", PORT)

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
					logger.Info("server is shutting down")
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
			logger.Error("error accepting connection", "details", err)
			break Accept
		}

		// No error happend, handle the connection in a separate goroutine
		wg.Add(1)
		go handleConnection(ctx, &wg, conn)
	}

	logger.Info("server stopped")
}

func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	// Close the connection when the function returns
	defer wg.Done()
	defer conn.Close()

	logger.Info("received connection", "from", conn.RemoteAddr().String())

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
				logger.Error("error reading from connection", "details", err)
				break Messages
			}

			// Print message to the console
			message = strings.TrimSpace(message)
			if message == "" {
				continue
			}
			logger.Info("received message", "message", message)

			// Write the message back to the connection
			_, err = conn.Write([]byte(message + "\n"))
			if err != nil {
				// Error writing to connection, log the error and stop
				// handling messages
				logger.Error("error writing to connection", "details", err)
				break Messages
			}
		}
	}

	logger.Info("connection closed")
}

func IsTimeout(err error) bool {
	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		return true
	}
	return false
}

func CreateSigintQuitChannel() <-chan os.Signal {
	// Create a channel to listen for SIGINT and SIGTERM signals.
	// When a signal is received, we signal the quit channel to
	// stop the server.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}

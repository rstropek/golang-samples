package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	customersvc "github.com/rstropek/golang-samples/go-kit"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":4000", "HTTP listen address")
	)
	flag.Parse()

	// Create logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create customer service and surround it with logging middleware
	var s customersvc.CustomerService
	{
		s = customersvc.NewCustomerRepository()
		s = customersvc.CustomerLoggingMiddleware(logger)(s)
	}

	// Create HTTP transport for customer service
	var h http.Handler
	{
		h = customersvc.MakeCustomerHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

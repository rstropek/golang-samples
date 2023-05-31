package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const (
	minBytes = 1
	maxBytes = 5 * 1024 * 1024 // 5MB
)

type App struct {
	client appinsights.TelemetryClient
}

func main() {
	app := App{}

	telemetryConfig := appinsights.NewTelemetryConfiguration(os.Getenv("INSTRUMENTATION_KEY"))
	telemetryConfig.MaxBatchSize = 8192
	telemetryConfig.MaxBatchInterval = 5 * time.Second
	app.client = appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	region := os.Getenv("REGION")
	if region == "" {
		region = "unknown"
	}
	app.client.Context().Tags.Cloud().SetRole(region)

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/singlebyte", app.requestTimingMiddleware(app.singleByteHandler))
	router.HandlerFunc(http.MethodPost, "/multibytes", app.requestTimingMiddleware(app.multiBytesHandler))
	router.HandlerFunc(http.MethodPost, "/proxy/:name", app.requestTimingMiddleware(app.proxyHandler))
	router.HandlerFunc(http.MethodPost, "/flush", app.requestTimingMiddleware(app.flushHandler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(code int) {
	rwi.statusCode = code
	rwi.ResponseWriter.WriteHeader(code)
}

func (app *App) requestTimingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rwi := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next(rwi, r)

		elapsed := time.Since(start)
		request := appinsights.NewRequestTelemetry(r.Method, r.RequestURI, elapsed, strconv.Itoa(rwi.statusCode))
		app.client.Track(request)
	}
}

func (app *App) singleByteHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 1)
	b[0] = 42
	w.Write(b)
}

func (app *App) multiBytesHandler(w http.ResponseWriter, r *http.Request) {
	lenStr := r.URL.Query().Get("len")
	length, err := strconv.Atoi(lenStr)
	if err != nil || length < minBytes || length > maxBytes {
		app.client.TrackException(fmt.Errorf("invalid length %s", lenStr))
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = 42
	}

	io.WriteString(w, string(b))
}

type ProxyRequest struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func (app *App) proxyHandler(w http.ResponseWriter, r *http.Request) {
	var req ProxyRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		app.client.TrackException(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	start := time.Now()
	dependency := appinsights.NewRemoteDependencyTelemetry(req.Name, "External", req.URL, true)
	resp, err := http.Post(req.URL, "text/plain", nil)
	if err != nil {
		app.dependencyError(err, w, dependency, start, resp, "Error making proxy request")
		return
	}
	defer resp.Body.Close()

	proxyBody, err := io.ReadAll(resp.Body)
	if err != nil {
		app.dependencyError(err, w, dependency, start, resp, "Error reading proxy response")
		return
	}
	dependency.Duration = time.Since(start)
	app.client.Track(dependency)

	w.Write(proxyBody)
}

func (app *App) flushHandler(w http.ResponseWriter, r *http.Request) {
	app.client.Channel().Flush()
	w.WriteHeader(http.StatusOK)
}

func (app *App) dependencyError(err error, w http.ResponseWriter, dependency *appinsights.RemoteDependencyTelemetry, start time.Time, resp *http.Response, errorDescription string) {
	app.client.TrackException(err)
	http.Error(w, errorDescription, http.StatusInternalServerError)
	dependency.Duration = time.Since(start)
	dependency.Success = false
	dependency.ResultCode = strconv.Itoa(resp.StatusCode)
	app.client.Track(dependency)
}

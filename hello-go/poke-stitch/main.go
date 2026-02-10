package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	gim "github.com/ozankasikci/go-image-merge"
)

// main wires together process lifecycle concerns:
// command-line parsing, dependency construction, HTTP server startup, and
// graceful shutdown. In Go, `main` commonly acts as explicit composition root
// instead of relying on framework magic.
func main() {
	// `flag` is Go's standard library option parser. It mutates package-level
	// state and therefore must be called before values are consumed.
	port := flag.Uint("p", 8080, "port to listen on")
	flag.Parse()

	// `slog` is the structured logging package introduced in the standard
	// library. We emit key/value fields to make machine processing easy.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// A shared HTTP client is preferred in Go:
	// - it reuses TCP connections via internal transports
	// - it centralizes timeout policy
	// Creating clients per request is an anti-pattern.
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Constructor returns `http.Handler` interface instead of concrete type.
	// Exposing behavior through interfaces is a common Go design style.
	handler := newStitchHandler(httpClient)

	// `http.ServeMux` in current Go versions supports method-aware patterns like
	// "GET /path". This avoids custom method checks inside handlers.
	mux := http.NewServeMux()
	mux.Handle("GET /poke-stitch", handler)

	address := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    address,
		Handler: mux,
		// Protect against slowloris-style attacks by bounding header read time.
		ReadHeaderTimeout: 5 * time.Second,
	}

	// `signal.NotifyContext` derives a context canceled by OS signals.
	// This avoids manual signal channels and integrates naturally with Go's
	// context-driven cancellation model.
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// `stop` unregisters signal handlers and should always be called.
	defer stop()

	// Run HTTP serving in a goroutine so main goroutine can block on shutdown.
	// Goroutines are lightweight user-space threads managed by the runtime.
	go func() {
		logger.Info("starting server", "address", address)
		// ListenAndServe normally returns `http.ErrServerClosed` on graceful
		// shutdown. Treat only other errors as fatal.
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			// Trigger shutdown path if startup/runtime unexpectedly fails.
			stop()
		}
	}()

	// Block until signal arrives or `stop()` is called due to server failure.
	<-shutdownCtx.Done()

	// Always bound graceful shutdown time; otherwise buggy handlers may keep
	// the process alive indefinitely.
	gracefulTimeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("shutting down server")
	// Shutdown stops accepting new connections and waits for active requests.
	if err := server.Shutdown(gracefulTimeoutCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
}

// pokemon matches only the PokeAPI fields this sample uses.
// Go's JSON decoder ignores unknown fields by default.
type pokemon struct {
	Sprites pokemonSprites `json:"sprites"`
}

// pokemonSprites mirrors selected sprite URLs from the API payload.
// Struct tags map Go field names to snake_case JSON keys.
type pokemonSprites struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

// stitchHandler holds request-scoped dependencies.
// In Go, handlers are often small structs with methods rather than closures.
type stitchHandler struct {
	client *http.Client
}

// newStitchHandler is a lightweight constructor. Returning `http.Handler`
// keeps callers decoupled from implementation details.
func newStitchHandler(client *http.Client) http.Handler {
	return &stitchHandler{client: client}
}

// ServeHTTP implements the `http.Handler` interface. Any type with this method
// can be registered in the router ("interface satisfaction by method set").
func (h *stitchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Query values are always strings; missing keys produce empty string.
	pokemonName := r.URL.Query().Get("pokemon")
	if pokemonName == "" {
		// `http.Error` writes status + plain text body in one call.
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Pass request context downstream so client disconnects and timeouts can
	// cancel outbound requests automatically.
	p, err := h.fetchPokemon(r.Context(), pokemonName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Retrieve available sprites concurrently.
	images, err := h.fetchSprites(r.Context(), p.Sprites.urls())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Defensive checks keep failure mode explicit for malformed upstream data.
	if len(images) == 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Preallocate slice capacity to avoid repeated allocations during append.
	grids := make([]*gim.Grid, 0, len(images))
	for _, img := range images {
		// Keep nil filtering local; merge library expects concrete images.
		if img != nil {
			grids = append(grids, &gim.Grid{Image: img})
		}
	}
	if len(grids) == 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Integer math trick: ceil(len/2) for two columns.
	rows := (len(grids) + 1) / 2
	rgba, err := gim.New(grids, 2, rows).Merge()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Header must be set before first body write.
	w.Header().Set("Content-Type", "image/png")
	// Encoding directly to response avoids temporary buffering.
	if err := png.Encode(w, rgba); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// fetchPokemon requests metadata JSON from PokeAPI and decodes it into `pokemon`.
// It returns value+error (Go's primary error handling convention).
func (h *stitchHandler) fetchPokemon(ctx context.Context, pokemonName string) (pokemon, error) {
	// Request is bound to caller context for cooperative cancellation.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://pokeapi.co/api/v2/pokemon/"+pokemonName, nil)
	if err != nil {
		return pokemon{}, err
	}

	// `Do` returns a response for any HTTP status; status validation is caller's
	// responsibility.
	res, err := h.client.Do(req)
	if err != nil {
		return pokemon{}, err
	}
	// Always close response body to return connection to pool.
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return pokemon{}, fmt.Errorf("unexpected status from pokeapi: %d", res.StatusCode)
	}

	var p pokemon
	// Streaming decoder avoids reading entire body into memory first.
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return pokemon{}, err
	}

	return p, nil
}

// fetchSprites downloads sprite images concurrently.
// We preserve output order by storing each result at its original index.
func (h *stitchHandler) fetchSprites(ctx context.Context, urls []string) ([]image.Image, error) {
	// Fixed-length output slice allows lock-free indexed writes from workers.
	images := make([]image.Image, len(urls))

	// WaitGroup tracks worker completion. `Add` must happen before goroutines
	// start to avoid race with `Wait`.
	var wg sync.WaitGroup

	// Capture first error across goroutines.
	var firstErr error
	var firstErrMu sync.Mutex

	// Closure serializes first-error assignment.
	recordErr := func(err error) {
		if err == nil {
			return
		}
		firstErrMu.Lock()
		defer firstErrMu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	wg.Add(len(urls))
	for i, url := range urls {
		// Pass loop variables as parameters. In Go, loop vars are reused each
		// iteration; explicit parameters avoid accidental capture bugs.
		go func(index int, spriteURL string) {
			defer wg.Done()
			img, err := h.getImage(ctx, spriteURL)
			if err != nil {
				recordErr(err)
				return
			}
			// Different goroutines write to distinct indices, which is safe.
			images[index] = img
		}(i, url)
	}

	// Join point: wait until all goroutines complete.
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return images, nil
}

// getImage performs one HTTP fetch + image decode operation.
func (h *stitchHandler) getImage(ctx context.Context, url string) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Explicit status handling keeps operational failures visible.
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from sprite endpoint: %d", res.StatusCode)
	}

	// `image.Decode` auto-detects registered formats (png/jpeg/gif/etc.).
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// urls returns only non-empty sprite URLs in deterministic order.
// Keeping filtering in one place simplifies caller logic.
func (s pokemonSprites) urls() []string {
	all := []string{
		s.BackDefault,
		s.BackFemale,
		s.BackShiny,
		s.BackShinyFemale,
		s.FrontDefault,
		s.FrontFemale,
		s.FrontShiny,
		s.FrontShinyFemale,
	}
	result := make([]string, 0, len(all))
	for _, u := range all {
		if u != "" {
			result = append(result, u)
		}
	}
	return result
}

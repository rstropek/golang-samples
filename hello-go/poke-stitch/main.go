package main

import (
    // Importing necessary packages
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sync"

	gim "github.com/ozankasikci/go-image-merge" // Importing a third-party package for image manipulation
)

// main is the entry point of the program.
func main() {
    // Setup for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

    // Command line flag parsing
	port := flag.Uint("p", 8080, "port to listen on")
	flag.Parse()

    // HTTP server setup
	mux := http.NewServeMux()
	mux.HandleFunc("/poke-stitch", stitch) // Registering the handler function

	address := fmt.Sprintf(":%d", *port)
	server := &http.Server{Addr: address, Handler: mux}
	go func() {
        // Starting the server in a goroutine
		server.ListenAndServe()
	}()

	<-stop // Waiting for interrupt signal
	fmt.Println("Stopping the server")
	server.Shutdown(context.Background()) // Graceful shutdown
	fmt.Println("We are done")
}

// pokemon struct to unmarshal JSON data from PokeAPI.
type pokemon struct {
	Sprites pokemonSprites `json:"sprites"`
}

// pokemonSprites struct to map the JSON sprite data.
type pokemonSprites struct {
	// Fields for different sprite URLs
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

// stitch is an HTTP handler function for the /poke-stitch endpoint.
func stitch(w http.ResponseWriter, r *http.Request) {
    // Extracting query parameter
	poke := r.URL.Query().Get("pokemon")
	if len(poke) == 0 {
		w.WriteHeader(http.StatusBadRequest) // Send bad request status if no pokemon specified
		return
	}

    // Fetching data from PokeAPI
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + poke)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    // Unmarshaling JSON data into pokemon struct
	var p pokemon
	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    // Setting up a channel and WaitGroup for concurrent image fetching
	images := make(chan image.Image, 8)
	wg := &sync.WaitGroup{}
	wg.Add(8)

    // Fetching each sprite image concurrently
	go getImage(p.Sprites.BackDefault, images, wg)
	go getImage(p.Sprites.BackFemale, images, wg)
	go getImage(p.Sprites.BackShiny, images, wg)
	go getImage(p.Sprites.BackShinyFemale, images, wg)
	go getImage(p.Sprites.FrontDefault, images, wg)
	go getImage(p.Sprites.FrontFemale, images, wg)
	go getImage(p.Sprites.FrontShiny, images, wg)
	go getImage(p.Sprites.FrontShinyFemale, images, wg)
	
	wg.Wait() // Waiting for all goroutines to finish
	close(images) // Closing the channel
	
    // Preparing images for merging
	grids := make([]*gim.Grid, 0)
	for img := range images {
		grids = append(grids, &gim.Grid{Image: img})
	}

    // Merging images
	rgba, err := gim.New(grids, 2, int(math.Ceil(float64(len(grids))/float64(2)))).Merge()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    // Encoding the merged image to PNG
	b := new(bytes.Buffer)
	wr := bufio.NewWriter(b)
	err = png.Encode(wr, rgba)
	wr.Flush()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    // Setting response headers and sending the image
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}

// getImage fetches an image from a URL and sends it to the provided channel.
func getImage(url string, out chan<- image.Image, wg *sync.WaitGroup) {
	defer wg.Done() // Marking this function as done in the WaitGroup upon return
	if len(url) == 0 {
		return // If the URL is empty, return immediately
	}

    // Fetching the image
	res, err := http.Get(url)
	if err != nil {
		return
	}

    // Decoding the image
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return
	}

    // Sending the image to the channel
	out <- img
}

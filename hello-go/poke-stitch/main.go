package main

import (
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

	gim "github.com/ozankasikci/go-image-merge"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	port := flag.Uint("p", 8080, "port to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/poke-stitch", stitch)

	address := fmt.Sprintf(":%d", *port)
	server := &http.Server{Addr: address, Handler: mux}
	go func() {
		server.ListenAndServe()
	}()

	<-stop
	fmt.Println("Stopping the server")
	server.Shutdown(context.Background())
	fmt.Println("We are done")
}

type pokemon struct {
	Sprites pokemonSprites `json:"sprites"`
}

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

func stitch(w http.ResponseWriter, r *http.Request) {
	poke := r.URL.Query().Get("pokemon")
	if len(poke) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	var p pokemon
	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	images := make(chan image.Image, 8)
	wg := &sync.WaitGroup{}
	wg.Add(8)

	go getImage(p.Sprites.BackDefault, images, wg)
	go getImage(p.Sprites.BackFemale, images, wg)
	go getImage(p.Sprites.BackShiny, images, wg)
	go getImage(p.Sprites.BackShinyFemale, images, wg)
	go getImage(p.Sprites.FrontDefault, images, wg)
	go getImage(p.Sprites.FrontFemale, images, wg)
	go getImage(p.Sprites.FrontShiny, images, wg)
	go getImage(p.Sprites.FrontShinyFemale, images, wg)
	
	wg.Wait()
	close(images)
	
	grids := make([]*gim.Grid, 0)
	for img := range images {
		grids = append(grids, &gim.Grid{Image: img})
	}

	rgba, err := gim.New(grids, 2, int(math.Ceil(float64(len(grids))/float64(2)))).Merge()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b := new(bytes.Buffer)
	wr := bufio.NewWriter(b)
	err = png.Encode(wr, rgba)
	wr.Flush()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}

func getImage(url string, out chan<- image.Image, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(url) == 0 {
		return
	}

	res, err := http.Get(url)
	if err != nil {
		return
	}

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return
	}

	out <- img
}

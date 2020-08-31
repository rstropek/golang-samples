package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func runInParallel(body func(chan<- uint, *sync.WaitGroup), goroutines int, agg func(uint, uint) uint, result chan<- uint) {
	results := make(chan uint, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go body(results, &wg)
	}
	wg.Wait()
	close(results)

	aggResult := uint(0)
	for item := range results {
		aggResult = agg(aggResult, item)
	}

	result <- aggResult
}

func main2() {
	cpus := runtime.NumCPU()
	// runtime.GOMAXPROCS(cpus)

	resultChan := make(chan uint)

	const ITERATIONS = 10000000

	go runInParallel(func(result chan<- uint, wg *sync.WaitGroup) {
		inside := uint(0)
		for i := 0; i < ITERATIONS; i++ {
			a := rand.Float32()
			b := rand.Float32()
			c := a*a + b*b
			if c <= float32(1) {
				inside++
			}
		}
		result <- inside
		wg.Done()
	}, cpus, func(a uint, b uint) uint { return a + b }, resultChan)

	result := <-resultChan
	pi := float32(result) / float32(cpus*ITERATIONS) * float32(4)
	fmt.Printf("%.6f", pi)
}

func main3() {
	cpus := runtime.NumCPU()
	// runtime.GOMAXPROCS(cpus)

	const ITERATIONS = 10000000

	results := make([]uint, cpus)

	var wg sync.WaitGroup
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		currentGoroutine := i
		go func() {
			inside := uint(0)
			r := rand.New(rand.NewSource(int64(currentGoroutine) * time.Now().UTC().UnixNano()))
			for i := 0; i < ITERATIONS; i++ {
				a := r.Float32()
				b := r.Float32()
				c := a*a + b*b
				if c <= float32(1) {
					inside++
				}
			}
			results[currentGoroutine] = inside
			wg.Done()
		}()
	}
	wg.Wait()

	aggResult := uint(0)
	for _, item := range results {
		aggResult = aggResult + item
	}

	pi := float32(aggResult) / float32(cpus*ITERATIONS) * float32(4)
	fmt.Printf("%.6f", pi)
}

func main() {
	const ITERATIONS = 10000000

	// Create a new random number generate. Attention! r is not thread safe.
	// You need to call `rand.New` in each goroutine.
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	inside := uint(0)
	for i := 0; i < ITERATIONS; i++ {
		a := r.Float32()
		b := r.Float32()
		c := a*a + b*b
		if c <= float32(1) {
			inside++
		}
	}

	pi := float32(inside) / float32(ITERATIONS) * float32(4)
	fmt.Printf("%.6f", pi)
}

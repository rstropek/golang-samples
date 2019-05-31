package main

import (
	"fmt"
	"time"

	ba "github.com/golang-collections/go-datastructures/bitarray"
	qpbas "github.com/rstropek/golang-samples/queens-problem/queens-problem-bitarray-solver"
)

func main() {
	fmt.Print("Starting calculation\n")
	solutions := findAllSolutions()
	// for _, s := range solutions {
	// 	qpbas.Print(s, os.Stdout)
	// }

	fmt.Printf("Found %d solutions\n", len(solutions))
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func findAllSolutions() []ba.BitArray {
	defer elapsed("Finding all solutions for queens problem")()
	return qpbas.FindSolution(ba.NewBitArray(8*8), 0, make([]ba.BitArray, 0))
}

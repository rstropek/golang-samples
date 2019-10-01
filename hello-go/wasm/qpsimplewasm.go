package main

import (
	"fmt"
	"os"

	qpbas "github.com/rstropek/golang-samples/hello-go/wasm/queensproblembitarraysolver"
)

func main() {
	var sl byte
	sl = 8
	fmt.Printf("Solving n queens problem for n=%d...\n", sl)

	result := qpbas.FindSolutions(sl)

	// Print first two solutions
	for ix, s := range result.Solutions {
		qpbas.Print(s, sl, os.Stdout)
		if ix == 1 {
			break
		}
	}

	fmt.Printf("Finding %d solutions took %v\n", len(result.Solutions), result.CalculationTime)
}

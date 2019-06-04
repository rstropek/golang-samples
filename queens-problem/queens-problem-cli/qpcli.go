package main

import (
	"flag"
	"fmt"
	"os"

	qpbas "github.com/rstropek/golang-samples/queens-problem/queens-problem-bitarray-solver"
)

func main() {
	sideLength := flag.Uint("sl", 8, "Side length of the chess board")
	printSolutions := flag.Bool("p", false, "Indicating whether solutions should be printed to Stdout")
	flag.Parse()

	sl := (byte)(*sideLength)
	fmt.Printf("Solving n queens problem for n=%d...\n", sl)

	result := qpbas.FindSolutions(sl)

	if *printSolutions {
		// Print all solutions
		for _, s := range result.Solutions {
			qpbas.Print(s, sl, os.Stdout)
		}
	}

	fmt.Printf("Finding %d solutions took %v\n", len(result.Solutions), result.CalculationTime)
}

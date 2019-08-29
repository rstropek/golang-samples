package main

import (
	"errors"
	"fmt"
)

// The following method returns the result AND an error. The error is nil if
// everything is ok.

func div(x int, y int) (int, error) {
	if y == 0 {
		return -1, errors.New("Sorry, division by zero is not supported")
	}

	return x / y, nil
}

func main() {
	// Here we declare-and-assign the result and the error variable.
	result, err := div(42, 0)
	if err != nil {
		fmt.Printf("Ups, something bad happened: %s\n", err)
		return
	}

	fmt.Printf("The result is %d\n", result)
}

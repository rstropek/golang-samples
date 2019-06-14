package calculator

import (
	"fmt"
)

// Note uppercase functions names -> exported

// Add adds two given numbers
func Add(x int, y int) int {
	return x + y
}

// Sub subtracts two given numbers
func Sub(x int, y int) int {
	return x - y
}

// Div divides two given numbers
func Div(x int, y int) (int, error) {
	if y != 0 {
		return x / y, nil
	}

	// Note multiple return values. By convention, last one is error
	// (see also https://golang.org/pkg/errors/).
	return 0, fmt.Errorf("y must not be zero")
}

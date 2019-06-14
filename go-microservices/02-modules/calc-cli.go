package main

import (
	"fmt"
	"os"
	"strconv"

	// Note path relative to GOPATH here
	"github.com/rstropek/golang-samples/go-microservices/02-modules/calculator"
)

func main() {
	// We manually handle command line args here. This is for demo purposes only. For more
	// complex CLIs, use a package like https://golang.org/pkg/flag/.

	// Note creating a slice (https://gobyexample.com/slices)
	// Also note declare-and-assign syntax
	args := os.Args[1:]

	// Note usage of the builtin function `len`
	// (see also https://golang.org/pkg/builtin/#len)
	if len(args) == 3 {
		// Note error handline
		x, convErr := strconv.Atoi(args[0])
		handleArgumentError(convErr)

		y, convErr := strconv.Atoi(args[2])
		handleArgumentError(convErr)

		var result int
		switch args[1] {
		case "+":
			result = calculator.Add(x, y)
		case "-":
			result = calculator.Sub(x, y)
		case "/":
			var err error
			result, err = calculator.Div(x, y)
			if err != nil {
				panic(fmt.Errorf("Cannot divide %d by %d", x, y))
			}
		}

		// For demo: Make mistake in format string regarding data types -> warning
		fmt.Printf("The result of %d %s %d is %d", x, args[1], y, result)
	} else {
		panic("Missing parameters")
	}
}

func handleArgumentError(err error) {
	if err != nil {
		// Note exiting goroutine (see also https://programming.guide/go/panic-explained.html)
		panic("Invalid argument")
	}
}

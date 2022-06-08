package main

import (
	"fmt"

	"github.com/mbndr/figlet4go"
	"rsc.io/quote"
)

func main() {
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(quote.Hello())

	fmt.Print(renderStr)
}

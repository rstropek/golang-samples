package main

import (
	"syscall/js"

	qpbas "github.com/rstropek/golang-samples/queens-problem/queens-problem-bitarray-solver"
)

var c chan bool

func findSolutions(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.Undefined()
	}

	sl := (byte)(args[0].Int())
	result := qpbas.FindSolutions(sl)

	return len(result.Solutions)
}

func main() {
	document := js.Global().Get("document")
	p := document.Call("createElement", "h1")
	p.Set("innerHTML", "Queens Problem")
	document.Get("body").Call("appendChild", p)

	p = document.Call("createElement", "p")
	p.Set("innerHTML", "Run the solver from console with `findSolutions(8)`")
	document.Get("body").Call("appendChild", p)

	c = make(chan bool)
	js.Global().Set("findSolutions", js.FuncOf(findSolutions))
	<-c
}

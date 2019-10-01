package main

// #include <stdio.h>
// #include <stdlib.h>
// #include "qps.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	numberOfSolutions := C.calculateNumberOfSolutions(8)
	fmt.Printf("We have found %d solutions\n", numberOfSolutions)

	greeting := C.CString("Hello World")
	defer C.free(unsafe.Pointer(greeting))
	C.puts(greeting)
}

package main

// How to compile C-code
// gcc -c -Wall -Werror -fpic qps.c
// gcc -shared -o libqps.so qps.o
// rm qps.o
// echo $LD_LIBRARY_PATH
// export LD_LIBRARY_PATH=$(pwd)

// #cgo LDFLAGS: -L. -lqps
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

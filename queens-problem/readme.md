# Queens Problem Solver Example

## Introduction

This set of examples can be used to demonstrate the use of [Go](https://golang.org/) in WebAssembly (WASM).

## Projects

* [*queens-problem-bitarray-solver*](queens-problem-bitarray-solver) contains the business logic incl. unit tests for solving the *n* queens problem.
* [*queens-problem-cli*](queens-problem-cli) contains a demo CLI using the solver.
* [*queens-problem-web*](queens-problem-web) contains a demo web API using the solver. This sample also contains a *Dockerfile* for creating a Docker image.
* [*queens-problem-simple-wasm*](queens-problem-simple-wasm) contains the solver in WebAssembly. It also shows a variant with [TinyGo](https://tinygo.org)
* [*queens-problem-js-solver*](queens-problem-js-solver) and [*queens-problem-cpp-solver*](queens-problem-cpp-solver) contain translations in JavaScript and C++. Use them to compare performance and size.
* [*queens-problem-wasm*](queens-problem-wasm) contains a variant of the solver in WebAssembly that demonstrates JavaScript interop.

## Demo Script

### Why I Wanted to Become a Gopher

![Gopher](https://blog.golang.org/gopher/header.jpg)

* C-family language
* Influenced by Pascal/Modula/Oberon
* Reduced clutter and complexity
  * [No generics](https://golang.org/doc/faq#generics)
  * [No exceptions](https://golang.org/doc/faq#exceptions)
* No type hierarchy
  * Composition instead of inheritance
  * Dispatching of the methods via [interfaces](https://gobyexample.com/interfaces)
* [Goroutines](https://golang.org/doc/faq#goroutines)
  * Multiplex functions (aka coroutines) on a set of threads
  * Growing and shrinking threads
  * [Channels](https://gobyexample.com/channels) (=pipes) to communicate between goroutines
* Runtime library, but no virtual machine
  * Code is compiled AOT to machine code
  * Cross-platform build
* Fast and efficient
  * Build
  * Execute
  * Size (e.g. containers)
* Great for CLIs and Web APIs
  * JSON, ProtoBuf, [gobs](https://golang.org/pkg/encoding/gob/)

### Our Example

* [*n* Queens Problem](https://en.wikipedia.org/wiki/Eight_queens_puzzle)
* [Simple implementation](queens-problem-bitarray-solver) in Go using [`bitarray`](https://godoc.org/github.com/golang-collections/go-datastructures/bitarray). A multi-threaded solution is not a goal of this example.
  * Code walkthrough ([code](queens-problem-bitarray-solver/queens-problem-bitarray-solver.go) and [unit tests](queens-problem-bitarray-solver/queens-problem-bitarray-solver_test.go))
  * Run tests with `go test`

### Warmup: Go CLI

* [CLI](queens-problem-cli)
  * [Code walkthrough](queens-problem-cli/qpcli.go)
  * Build with `go build`
  * Run with `qpcli -p`

### WASM Baseline

* [JavaScript translation](queens-problem-js-solver)
* [C++/Emscripten translation](queens-problem-cpp-solver) (thanks to [@ArnoHu](https://twitter.com/arnohu) for his support)


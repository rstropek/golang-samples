# Queens Problem Solver Example

## Introduction

This set of examples can be used to demonstrate the use of Go in WebAssembly (WASM).

## Projects

* [*queens-problem-bitarray-solver*](queens-problem-bitarray-solver) contains the business logic incl. unit tests for solving the *n* queens problem.
* [*queens-problem-cli*](queens-problem-cli) contains a demo CLI using the solver.
* [*queens-problem-web*](queens-problem-web) contains a demo web API using the solver. This sample also contains a *Dockerfile* for creating a Docker image.
* [*queens-problem-simple-wasm*](queens-problem-simple-wasm) contains the solver in WebAssembly. It also shows a variant with [TinyGo](https://tinygo.org)
* [*queens-problem-js-solver*](queens-problem-js-solver) and [*queens-problem-cpp-solver*](queens-problem-cpp-solver) contain translations in JavaScript and C++. Use them to compare performance and size.
* [*queens-problem-wasm*](queens-problem-wasm) contains a variant of the solver in WebAssembly that demonstrates JavaScript interop.

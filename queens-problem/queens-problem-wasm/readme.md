# Queens Problem Solver WASM

## Introduction

This Go example runs a [*n* queens problem](https://en.wikipedia.org/wiki/Eight_queens_puzzle) solver using WebAssembly in the browser. The solver example can be found [here](../queens-problem-bitarray-solver). This example demonstrates JavaScript interop.

## How to Use

* Run `docker build -t qpwasm .` to create Docker image hosting the website. Run the web API with `docker run -d -p 8080:80 --name qpwasm qpwasm`.
* Open [http://localhost:8080](http://localhost:8080) and see output of Go app in the console window
* Run the solver from console with `findSolutions(8)`

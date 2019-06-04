# Queens Problem Solver WASM

## Introduction

This Go example runs a [*n* queens problem](https://en.wikipedia.org/wiki/Eight_queens_puzzle) solver using WebAssembly in the browser. The solver example can be found [here](../queens-problem-bitarray-solver).

## How to Use

* Run `go build -o qpsimplewasm.exe` to build an executable
* Run `docker build -t qpsimplewasm .` to create Docker image hosting the website. Run the web API with `docker run -d -p 8080:80 --name qpsimplewasm qpsimplewasm`.
* Open [http://localhost:8080](http://localhost:8080) and see output of Go app in the console window. Note content compression in network tab.
* Run `docker build -t qptinywasm -f Dockerfile.tiny .` to create Docker image with [TinyGo](https://tinygo.org) hosting the website. Run the web API with `docker run -d -p 8080:80 --name qptinywasm qptinywasm`.
* Open [http://localhost:8080](http://localhost:8080) and see output of Go app in the console window. Note content compression in network tab. Compare size and runtime with full Go version.

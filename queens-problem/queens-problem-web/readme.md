# Queens Problem Solver Web Server

## Introduction

This Go example runs a [*n* queens problem](https://en.wikipedia.org/wiki/Eight_queens_puzzle) solver in a web API. The solver example can be found [here](../queens-problem-bitarray-solver).

## How to Use

* Run `go build -o qpweb.exe` to build
* Run `docker build -t qpweb .` to create Docker image for web API. Run the web API with `docker run -d -p 8080:80 --name qpweb qpweb`.
* Execute sample request in [*demo.http*]

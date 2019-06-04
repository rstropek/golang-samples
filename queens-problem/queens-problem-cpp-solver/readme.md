# Queens Problem Solver WASM

## Introduction

This is a C++ translation of the [*n* queens problem](https://en.wikipedia.org/wiki/Eight_queens_puzzle) solver written in Go from [here](../queens-problem-bitarray-solver). Use it to compare performance and size of WASM with JavaScript.

## How to Use

* Run `docker run --rm -v C:\Users\r.stropek\go\src\github.com\rstropek\golang-samples\queens-problem\queens-problem-cpp-solver:/app -w /app trzeci/emscripten g++ qps.cpp -I/usr/include/x86_64-linux-gnu/c++/6/ -O3 -Wno-c++11-extensions -o qpcpp` to build an executable that can be run locally
* Run `docker build -t qpcpp .` to create Docker image hosting the website. Run the web API with `docker run -d -p 8080:80 --name qpcpp qpcpp`.
* Open [http://localhost:8080](http://localhost:8080) and see output of WASM app in the console window. Note content compression in network tab.

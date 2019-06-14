# Script

## Hello Go!

* Example: [01-hello-go](01-hello-go)
* Code walkthrough [hello.go](01-hello-go/hello.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go help`
  * `go run .`
  * `go build`
  * `set GOOS=linux`, `go build`

## Calculator

* Example: [02-modules](02-modules)
* Code walkthrough:
  * [calculator.go](02-modules/calculator/calculator.go)
  * Unit tests [calculator_test.go](02-modules/calculator/calculator_test.go)
  * CLI [calc-cli.go](02-modules/calc-cli.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go test`
  * `go run . 21 + 21`
  * `go run . 42 / 0`

## Simple Web Server

* Example: [03-simple-web-server](03-simple-web-server)
* Code walkthrough:
  * [web.go](03-simple-web-server/web.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go run .`
  * `go run . --help`
  * `go run . -p 8081`
* Run [demo requests](03-simple-web-server/demo.http)

## Advanced Web Server

* Example: [04-advanced-web-server](04-advanced-web-server)
* Code walkthrough:
  * [web.go](04-advanced-web-server/web.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go run .`, stop server with *Ctrl+c*

## Simple Web API

* Example: [05-simple-web-api](05-simple-web-api)
* Code walkthrough:
  * [web.go](05-simple-web-api/web.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go run .`
* Run [demo requests](05-simple-web-api/demo.http)

## Advanced Web API

* Example: [06-advanced-web-api](06-advanced-web-api)
* Code walkthrough:
  * [cart.go](06-advanced-web-api/cart.go)
  * [shopping-cart-api.go](06-advanced-web-api/shopping-cart-api.go)
* [`go` Command](https://golang.org/cmd/go/):
  * `go run .`
* Run [demo requests](06-advanced-web-api/demo.http)

## Docker Container

* Example: [05-simple-web-api](05-simple-web-api)
* Code walkthrough:
  * [Dockerfile](05-simple-web-api/Dockerfile)
* Commands:
  * `docker build -t go-web .`
  * `docker run -d --name go-web -p 8080:80 go-web`
  * `docker rm -f go-web`
  * Change `FROM alpine:latest` to `FROM scratch` and demo difference in size
  * `docker tag go-web rstropek/go-web`
  * `docker push rstropek/go-web`
  * Show GO web API running in *Azure Container Instance*

## Azure DevOps

* Example: [05-simple-web-api](05-simple-web-api)
* Code walkthrough:
  * [azure-pipelines.yml](05-simple-web-api/azure-pipelines.yml)
* Create *Azure DevOps Pipeline* and show building of GO app

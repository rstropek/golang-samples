# Web API - End to End

## Preparation

* Install Go
* Install latest version of Visual Studio Code
* Install Go extension for VSCode
* Install an API client of your choice (e.g. Postman, Thunder Client-extension for VSCode, REST client-extension for VSCode)
* Install [Hey tool](https://github.com/rakyll/hey) for load testing

## Create folders

```bash
mkdir -p bin cmd/api internal/data
touch Makefile
touch cmd/api/main.go
```

* *bin* directory will contain our compiled application binaries
* *cmd/api* directory will contain the application-specific code
* *internal* directory will contain helper packages used by our API (=code that isn't application-specific and can potentially be reused)

## First steps

* *main.go*: *0010-hello-world*
* `go run ./cmd/api`
* *main.go*: *0020-api-starter*
* *healthcheck.go*: *0030-api-starter-healthcheck*
* `go run ./cmd/api`
* Try request *Healthcheck*

## Routing

* Router used: [httprouter](https://github.com/julienschmidt/httprouter)
* `touch cmd/api/routes.go`
* `go get github.com/julienschmidt/httprouter`
* *routes.go*: *0040-basic-router*
* *main.go*: Replace servemux and creation of http server with *0050-http-server-with-router*
* `touch cmd/api/heroes.go`
* *heroes.go*: *0060-basic-hero-handlers*
* `go run ./cmd/api`
* Try requests *CreateHero*, *ShowHero*, *HeroNotFound*
* `touch cmd/api/helpers.go`
* *helpers.go*: *0070-get-id-param*
* *heroes.go*: Replace `showHeroHandler` with *0080-update-show-hero*
* Try requests *CreateHero*, *ShowHero*, *HeroNotFound*

## JSON

* *healthcheck.go*: Replace `healthcheckHandler` with *0090-healthcheck-json*
* `go run ./cmd/api`
* Try request *Healthcheck*
* *helpers.go*: Add `writeJSON` with *0100-writeJSON*
* *healthcheck.go*: Replace `healthcheckHandler` with *0110-healthcheck-writeJSON*

## Heroes Structures

* `touch internal/data/heroes.go`
* *internal/data/heroes.go*: *0120-hero*
* *cmd/api/heroes.go*: Replace `showHeroHandler` with *0130-show-hero*
* `go run ./cmd/api`
* Try request *ShowHero*

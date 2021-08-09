# Web API - End to End

## Preparation

* Install Go
* Install latest version of Visual Studio Code
* Install Go extension for VSCode
* Install an API client of your choice (e.g. Postman, Thunder Client-extension for VSCode, REST client-extension for VSCode)
* Install [Hey tool](https://github.com/rakyll/hey) for load testing

## Create folders

```bash
mkdir -p bin cmd/api internal/data internal/validator
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

## Sending JSON

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
* *helpers.go*: Replace `writeJSON` with *0140-writeJSON-envelope*
* *cmd/api/heroes.go*: Replace `showHeroHandler` with *0150-showHero-envelope*
* `go run ./cmd/api`
* Try requests *Healthcheck*, *ShowHero*
* `touch internal/data/runtime.go`
* *runtime.go*: *0170-MarshalJSON-CanFly*
* `go run ./cmd/api`
* Try request *ShowHero*

## Error Responses

* `touch cmd/api/errors.go`
* *errors.go*: *0180-errors*
* *healthcheck.go*: Replace error handling for `writeJSON` with *0190-error*
* *cmd/api/heroes.go*: Replace error handling for `writeJSON` with *0190-error*
* *routes.go*: Add error routes with *0200-routing-errors*
* `go run ./cmd/api`
* Try requests *RouteDoesNotExist*, *WrongHttpMethod*

## Parsing JSON

* *cmd/api/heroes.go*: Replace `createHeroHandler` with *0210-createHero-json-parse*
* `go run ./cmd/api`
* Try requests *CreateHero*
* Debug
  * Add *.vscode/launch.json*:

    ```json
    {
        "version": "0.2.0",
        "configurations": [
            {
                "name": "Launch Package",
                "type": "go",
                "request": "launch",
                "mode": "auto",
                "program": "${workspaceFolder}/cmd/api"
            }
        ]
    }
    ```

  * Set breakpoint in *cmd/api/heroes.go* and try debugging
* *helpers.go*: *0220-readJSON*
* *cmd/api/heroes.go*: Replace `err := json.NewDecoder...` with *0230-apply-readJSON*
* `go run ./cmd/api`
* Try requests *BadlyFormedJson1*, *BadlyFormedJson2*, *WrongCanFly*
* `touch internal/validator/validator.go`
* *validator.go*: *0240-validation-helpers*
* *internal/data/heroes.go*: Add `ValidateHero` with *0250-validate-hero*
* *cmd/api/heroes.go*: Add validation checks after `app.readJSON` with *0260-use-validations*
* `go run ./cmd/api`
* Try requests *MissingName*, *InvalidCoolness*, *DuplicateTags*

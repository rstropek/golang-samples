REM Initialize current folder as a module
go mod init github.com/rstropek/golang-samples/hello-go/modules

REM Look at go.mod
type go.mod

REM Add sample code to main.go:
REM Snippet: go, 01-module
go run .

REM Look at go.mod again
type go.mod

REM Show folder C:\Users\r.stropek\go\pkg\mod\github.com\mbndr
dir C:\Users\r.stropek\go\pkg\mod\github.com\mbndr

###

REM Build executable
go build
dir modules.exe
modules.exe

REM Cross-compile to Linux on Windows, run it on Linux
set GOARCH=amd64
set GOOS=linux
go build
bash
./modules
REM DON'T FORGET to reset env variables

###

REM Add sample Dockerfile main.go:
REM Snippets: go, 02-Dockerfile
REM           go, 03-dockerignore

REM Build image based on Alpine
docker build -t hello-go .
docker run -t --rm hello-go

REM Check image size
docker images hello-go

REM Switch to „FROM scratch“ (no RUN) and build/run again
REM Check image size
docker images hello-go

REM IN wasm FOLDER

REM Cross-compile to WASM
set GOOS=js
set GOARCH=wasm
go build -o qpsimplewasm.wasm
dir *.wasm
REM DON'T FORGET to reset env variables

REM In wasm/queensproblembitarraysolver folder...
cd queensproblembitarraysolver
go test
cd ..

REM In wasm folder...
go build -o qpsimplewasm.exe
qpsimplewasm.exe

docker build -t qpsimplewasm .
docker run -d -p 8080:80 --name qpsimplewasm qpsimplewasm
start http://localhost:8080
docker rm -f qpsimplewasm

REM In c-interop folder...
docker run --rm -v C:\Code\GitHub\golang-samples\hello-go\c-interop:/app -w /app golang go build .
bash -c ./app

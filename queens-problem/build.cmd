set GOOS=windows
set GOARCH=amd64
go build -o qpc.exe

set GOOS=js
set GOARCH=wasm
go build -o qpc.wasm
copy %GOROOT%misc\wasm\wasm_exec.js .

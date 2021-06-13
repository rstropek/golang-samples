# Use a multi-stage build
FROM golang:latest AS builder

# Compile Go into WASM
WORKDIR /app
COPY ./go.mod ./
COPY ./*.go ./
ENV CGO_ENABLED=0
ENV GOOS=js
ENV GOARCH=wasm
RUN go get \
    && go build -ldflags="-s -w" -a -o ./qpwasm.wasm .

FROM nginx:alpine

# Copy exe from build container
COPY --from=builder /app/qpwasm.wasm /usr/share/nginx/html
COPY --from=builder /usr/local/go/misc/wasm/wasm_exec.js /usr/share/nginx/html
COPY *.html /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

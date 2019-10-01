# Use a multi-stage build
FROM golang:latest AS builder

# Compile Go into exe
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags="-s -w" -a -o ./qpsimplewasm.wasm .

FROM nginx:alpine

# Copy exe from build container
COPY --from=builder /app/qpsimplewasm.wasm /usr/share/nginx/html
COPY --from=builder /usr/local/go/misc/wasm/wasm_exec.js /usr/share/nginx/html
COPY *.html /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

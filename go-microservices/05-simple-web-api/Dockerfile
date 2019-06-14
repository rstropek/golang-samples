# Use a multi-stage build
FROM golang:latest AS builder

# Install module for solving n queens problem
RUN go get \
    github.com/gorilla/mux

# Compile Go into exe
WORKDIR /app
COPY ./*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o ./web .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy exe from build container
COPY --from=builder /app/web ./
RUN chmod +x ./web

# Define port on which the container will listen
EXPOSE 8080

# Define start command
CMD ["./web"]

# Use a multi-stage build
FROM golang:latest AS builder

# Install module for solving n queens problem
RUN go get \
    github.com/rstropek/golang-samples/queens-problem/queens-problem-bitarray-solver \
    github.com/gorilla/mux \
    github.com/gorilla/handlers

# Compile Go into exe
WORKDIR /app
COPY ./*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o ./qpweb .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy exe from build container
COPY --from=builder /app/qpweb ./
RUN chmod +x ./qpweb

# Define port on which the container will listen
EXPOSE 80

# Define start command
CMD ["./qpweb", "-p", "80"]

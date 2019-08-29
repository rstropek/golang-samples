# Use a multi-stage build
FROM golang:latest AS builder

# Compile Go into exe
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o ./modules .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy exe from build container
COPY --from=builder /app/modules ./
RUN chmod +x ./modules

# Define start command
CMD ["./modules"]

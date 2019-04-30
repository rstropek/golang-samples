# Use a multi-stage build
FROM golang:latest AS builder

# Install module for turning markdown into HTML
RUN go get github.com/shurcooL/github_flavored_markdown

# Compile Go into exe
WORKDIR /app
COPY ./*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o ./mdweb .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy exe from build container
COPY --from=builder /app/mdweb ./
RUN chmod +x ./mdweb

# Create folder from which markdown content is read
RUN mkdir -p /usr/share/mdweb/content
ENV CONTENT=/usr/share/mdweb/content

# Define port on which the container will listen
ENV PORT=80
EXPOSE 80

# Define start command
CMD ["./mdweb"]

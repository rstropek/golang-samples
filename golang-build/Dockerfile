FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY demo ./
RUN chmod +x ./demo

CMD ["./demo", "-port", ":80"]

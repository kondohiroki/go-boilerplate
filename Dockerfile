# Build stage
FROM golang:1.20.2-alpine3.17 as builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

RUN go clean -modcache

# Run stage
FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]

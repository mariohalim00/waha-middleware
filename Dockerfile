FROM golang:1.22-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN ./goose -dir ./internal/database/migrations \
    postgres "host=localhost user=postgres password=postgres dbname=waha-middleware sslmode=disable" up || echo "Goose migration skipped or failed (safe for container build)"

RUN go build -o server ./cmd/server

FROM alpine:latest

# Install ca-certificates if needed by your app (e.g., HTTPS requests)
RUN apk update && apk upgrade && apk add --no-cache ca-certificates

WORKDIR /app

# Copy only required artifacts
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./server"]

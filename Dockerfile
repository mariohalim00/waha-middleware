FROM golang:1.23-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM golang:1.23-alpine

# Install ca-certificates if needed by your app (e.g., HTTPS requests)
RUN apk update && apk upgrade && apk add --no-cache ca-certificates

WORKDIR /app

# Copy only required artifacts
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
COPY --from=builder /app ./projects

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN ./goose -dir ./projects/internal/database/migrations \
    postgres "host=localhost user=postgres password=suntzu2025 dbname=waha-middleware sslmode=disable" up || echo "Goose migration skipped or failed (safe for container build)"

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./server"]

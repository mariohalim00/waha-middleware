FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --link --from=builder /app/server .
COPY --link --from=builder /app/.env .
EXPOSE 8080
CMD ["./server"]
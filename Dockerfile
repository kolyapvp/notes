# Build Stage
FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot ./cmd/bot

# Final Image
FROM debian:bullseye-slim
WORKDIR /root/

# Установить ca-certificates для работы HTTPS
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bot .
COPY --from=builder /app/.env .

CMD ["./bot"]
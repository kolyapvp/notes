# Build Stage
FROM golang:1.21-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install ca-certificates and build the app
RUN go build -o app ./cmd/bot

# Final Image
FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/.env .

# Install ca-certificates in the final image

CMD ["./app"]
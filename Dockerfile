# Используем официальный образ Go для сборки
FROM golang:1.21-bullseye AS builder
WORKDIR /app

# Копируем файлы и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходные файлы
COPY . .

# Собираем приложение
RUN go build -o app ./cmd/bot

# Используем минимальный образ для финального контейнера
FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/.env .

# Указываем команду запуска
CMD ["./app"]
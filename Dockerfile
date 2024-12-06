# Используем базовый образ Golang
FROM golang:1.21-bullseye AS builder

WORKDIR /app

# Копируем исходный код
COPY . .

# Обновляем и устанавливаем зависимости, включая обновление glib

# Загружаем зависимости проекта и собираем бинарник
RUN go mod tidy
RUN go build -o app ./cmd/bot

# Используем минимальный образ для запуска
FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/.env .

# Запускаем приложение
CMD ["./app"]
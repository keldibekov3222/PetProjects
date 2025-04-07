# Используем официальный образ Golang для сборки
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app


# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код в контейнер
COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=0 go build -o /app/main .
# Используем легковесный образ для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder-контейнера
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config/config.env /app/config/config.env
COPY --from=builder /app/db/migrations /app/db/migrations

# Даём права на выполнение
RUN chmod +x /app/main

# **ЯВНО УКАЗЫВАЕМ ЗАПУСК `/app/main`**
ENTRYPOINT ["/app/main"]

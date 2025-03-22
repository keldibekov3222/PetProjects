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
WORKDIR /root/

# Копируем бинарник из builder-контейнера
COPY --from=builder /app/main /root/main
COPY config/config.env config/config.env

# Даём права на выполнение
RUN chmod +x ./main

# **ЯВНО УКАЗЫВАЕМ ЗАПУСК `./main`**
ENTRYPOINT ["./main"]

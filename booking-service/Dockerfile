# Этап 1: Сборка
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем необходимые инструменты для сборки
RUN apk add --no-cache git ca-certificates

# Копируем только файлы зависимостей для лучшего кэширования
COPY go.mod go.sum ./

# Кэшируем модули Go для ускорения повторных сборок
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Копируем исходный код
COPY . .

# Собираем статически связанный бинарник с оптимизациями
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/

# Этап 2: Запуск
FROM alpine:3.19

WORKDIR /app

# Устанавливаем сертификаты для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .

# Переменные окружения (значения по умолчанию, можно переопределить в docker-compose)
ENV DB_HOST=booking-postgres
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=booking_db
ENV DB_PORT=5432
ENV PORT=8082
ENV LOG_LEVEL=info
ENV KAFKA_BROKER=kafka:9092
ENV CINEMA_SERVICE_URL=http://cinema-service:8081

# Открываем порт
EXPOSE 8082

CMD ["./main"]
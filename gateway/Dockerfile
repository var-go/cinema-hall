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
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Этап 2: Запуск
FROM alpine:3.19

WORKDIR /app

# Устанавливаем сертификаты для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .

# Переменные окружения (значения по умолчанию, можно переопределить в docker-compose)
ENV PORT=8085
ENV USER_SERVICE_URL=http://user-service:8080
ENV MOVIE_SERVICE_URL=http://movie-service:8083
ENV CINEMA_SERVICE_URL=http://cinema-service:8081
ENV BOOKING_SERVICE_URL=http://booking-service:8082
ENV JWT_SECRET=your-secret-key-change-in-production

# Открываем порт
EXPOSE 8085

CMD ["./main"]


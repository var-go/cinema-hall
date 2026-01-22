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
ENV PORT=8083
ENV DATABASE_URL=postgres://postgres:postgres@postgres:5432/movie_db?sslmode=disable
ENV LOG_LEVEL=info

# Открываем порт
EXPOSE 8083

CMD ["./main"]


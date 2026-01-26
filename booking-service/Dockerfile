FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

ENV DB_HOST=booking-postgres
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=booking_db
ENV DB_PORT=5432
ENV PORT=8082
ENV LOG_LEVEL=info
ENV KAFKA_BROKER=kafka:9092
ENV CINEMA_SERVICE_URL=http://cinema-service:8081

EXPOSE 8082

CMD ["./main"]
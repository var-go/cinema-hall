FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/app/

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

ENV PORT=8080
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASS=postgres
ENV DB_NAME=user_db
ENV DB_SSLMODE=disable
ENV KAFKA_BROKER=kafka:9092
ENV JWT_SECRET=your-secret-key-change-in-production
ENV BOOKING_SERVICE_URL=http://booking-service:8082

EXPOSE 8080

CMD ["./main"]

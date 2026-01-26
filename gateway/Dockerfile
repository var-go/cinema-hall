FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

ENV PORT=8085
ENV USER_SERVICE_URL=http://user-service:8080
ENV MOVIE_SERVICE_URL=http://movie-service:8083
ENV CINEMA_SERVICE_URL=http://cinema-service:8081
ENV BOOKING_SERVICE_URL=http://booking-service:8082
ENV JWT_SECRET=your-secret-key-change-in-production

EXPOSE 8085

CMD ["./main"]


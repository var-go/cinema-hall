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

ENV PORT=8083
ENV DATABASE_URL=postgres://postgres:postgres@postgres:5432/movie_db?sslmode=disable
ENV LOG_LEVEL=info

EXPOSE 8083

CMD ["./main"]


APP_NAME=cinema-service
CMD_PATH=./cmd/main.go
DOCKER_IMAGE=cinema-service

.PHONY: build run test clean docker-build docker-run

run:
	go run $(CMD_PATH)

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

test:
	go test ./...

clean:
	rm -rf bin

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --rm -p 8082:8082 $(DOCKER_IMAGE)

.PHONY: run build test fmt vet lint tidy clean dev seed

GO           ?= go
BINARY       ?= cmd
CMD_MAIN     := ./cmd/main.go

run: ## Запуск основного приложения (HTTP-сервер)
	$(GO) run $(CMD_MAIN)

dev: ## Запуск в режиме разработки с hot reload (air)
	air -c .air.toml

build: ## Сборка бинарника приложения
	mkdir -p tmp
	$(GO) build -o tmp/$(BINARY) $(CMD_MAIN)

test: ## Запуск всех тестов
	$(GO) test -v ./...

test-cover:
	$(GO) test -cover ./...	

fmt: ## Форматирование кода
	$(GO) fmt ./...

cover-html:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

vet: ## Статический анализ кода
	$(GO) vet ./...

lint: ## Линтинг кода с помощью golangci-lint
	golangci-lint run

tidy: ## Обновление зависимостей (go.mod / go.sum)
	$(GO) mod tidy

clean: ## Удаление собранных бинарников
	rm -rf tmp

.PHONY: build run migrate test test-integration lint docker-up docker-down

# Сборка и запуск контейнеров
build:
	docker-compose up --build

# Запуск контейнеров
run:
	docker-compose up

# Запуск миграций
migrate:
	docker-compose run avito-shop-service /build -migrate

# Запуск всех тестов
test:
	go test -v ./... -cover

# Запуск интеграционных/E2E тестов
test-integration:
	go test -v -tags=integration ./integration

# Запуск линтера
lint:
	golangci-lint run

# Запуск контейнеров в фоновом режиме
docker-up:
	docker-compose up -d

# Остановка и удаление контейнеров
docker-down:
	docker-compose down

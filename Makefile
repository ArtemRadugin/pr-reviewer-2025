APP_NAME=pr-service

build:
	go build -o $(APP_NAME) ./cmd/main.go

run:
	go run ./cmd/main.go

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

migrate-up:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:qwerty@localhost:5432/prdb?sslmode=disable" up

migrate-down:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:qwerty@localhost:5432/prdb?sslmode=disable" down

.PHONY: build run docker-build docker-up docker-down migrate-up migrate-down

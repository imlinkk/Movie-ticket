.PHONY: run test build docker-up docker-down swagger clean

run:
	go run cmd/api/main.go

test:
	go test -v ./tests/...

build:
	go build -o bin/api ./cmd/api

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

clean:
	rm -rf bin/ coverage.out

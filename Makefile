run:
	@go run cmd/main.go

dev:
	@air

docker:
	@docker-compose --env-file .env up --build

swagger:
	@swag fmt
	@swag init -g ./cmd/main.go -o ./docs

up:
	@go run cmd/migrate/main.go up

down:
	@go run cmd/migrate/main.go down

test:
	@go test -v ./internal...

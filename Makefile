dev:
	@air

docker:
	@docker-compose --env-file .env up --build

swagger:
	@swag fmt
	@swag init -g ./cmd/main.go -o ./docs

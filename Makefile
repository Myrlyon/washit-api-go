dev:
	@air

docker:
	@docker-compose --env-file .env up --build
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


MODEL ?= user
test-file:
	@go test -v ./internal/$(MODEL)/service
	@go test -v ./internal/$(MODEL)/repository
# make test-file MODEL=order

MOCK_NAME ?= IUserService
MOCK_DIR ?= /user/service
make-mock:
	@mockery --name $(MOCK_NAME) --dir ./internal$(MOCK_DIR) --output ./internal$(MOCK_DIR)/mock
# make make-mock MOCK_NAME=IOtherService MOCK_DIR=/order/service

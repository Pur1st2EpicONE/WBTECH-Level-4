.PHONY: all L4.3 clean test lint

-include .env.example .env

all: L4.3

L4.3:
	@go build -o calendar ./cmd/calendar/main.go
	@./calendar

local:
	if [ ! -f .env ]; then cat .env.example > .env; fi 
	if [ ! -f config.yaml ]; then cp ./configs/config.dev.yaml ./config.yaml; fi 
	if [ ! -f docker-compose.yaml ]; then cp ./deployments/docker-compose.dev.yaml ./docker-compose.yaml; fi
	COMPOSE_BAKE=true docker compose up -d postgres
	until docker exec postgres pg_isready -U ${DB_USER} > /dev/null 2>&1; do sleep 1; done
	bash -c 'trap "exit 0" INT; go run ./cmd/calendar/main.go'

down:
	docker compose down 2>/dev/null || true 
	rm -f Dockerfile docker-compose.yaml config.yaml

clean:
	@rm -f ./calendar
	@rm -rf ./logs

test: 
	@go test ./internal/handler/v1 -cover
	@go test ./internal/service/impl -cover
	@go test ./internal/repository/memory -cover 

lint:
	golangci-lint run ./...

.env:
	@:

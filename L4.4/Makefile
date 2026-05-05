.PHONY: all up monitoring-up monitoring-down reset lint help
.POSIX:
.SILENT:

all: up

up:
	@go build -o profiler ./cmd/profiler/main.go
	@./profiler

monitoring-up:
	@COMPOSE_BAKE=true docker compose -f monitoring/docker-compose.yaml up -d

monitoring-down:
	@COMPOSE_BAKE=true docker compose -f monitoring/docker-compose.yaml down

reset:
	@rm -f ./profiler
	@docker compose -f monitoring/docker-compose.yaml down -v
	@rm -rf ./logs cpu.pprof

lint:
	@golangci-lint run ./...

test: 
	@go test ./internal/handler -cover

help:
	@echo " ———————————————————————————————————————————————————————————— "
	@echo "| up             | Build and run the profiler                |"
	@echo "| monitoring-up  | Start Prometheus and Grafana in Docker    |"
	@echo "| monitoring-down| Stop monitoring stack                     |"
	@echo "| reset          | Remove binary, containers and logs        |"
	@echo "| lint           | Run golangci-lint                         |"
	@echo "| test           | Run unit tests                            |"
	@echo " ———————————————————————————————————————————————————————————— "

.DEFAULT:
	@echo " No rule to make target '$@'. Available make targets:"
	@$(MAKE) --no-print-directory help
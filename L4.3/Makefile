.PHONY: all L2.18 clean test lint

all: L2.18

L2.18:
	@go build -o calendar ./cmd/calendar/main.go
	@./calendar

clean:
	@rm -f ./calendar
	@rm -rf ./logs

test: 
	@go test ./internal/handler/v1 -cover
	@go test ./internal/service/impl -cover
	@go test ./internal/repository/memory -cover 

lint:
	golangci-lint run ./...

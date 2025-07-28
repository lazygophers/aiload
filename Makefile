.PHONY: run
run:
	@go run ./cmd/main.go

.PHONY: test
test:
	@go test -v ./...

.PHONY: air
air:
	@air

.PHONY: cover
cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
tidy:
	@go mod tidy
run: tidy
	@go run cmd/main.go
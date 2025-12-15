unit-test:
	@echo ">> Running all tests in all packages"
	go test ./... -v -cover


coverage-test:
    go test ./... -coverprofile=coverage
    go tool cover -html=coverage

run: 
	@echo ">> Running the application"
	go run ./cmd/api
	
build:
	@echo ">> Building the application"
	go build -o bin/app ./cmd/api
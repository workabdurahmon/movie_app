.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Build the application
build:
	go build -o bin/movie_app cmd/api/main.go

# Run the application locally
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t movie_app .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env movie_app

# Start development environment with Docker Compose
docker-compose-up:
	docker-compose up --build

# Stop development environment
docker-compose-down:
	docker-compose down

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Run formatter
fmt:
	go fmt ./... 
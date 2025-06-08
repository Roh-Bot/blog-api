.PHONY: all deps build run debug image container stop-container compose db-migrate-up db-migrate-down swag-gen

image_version=1.0.0
image_name=go-clean:$(image_version)
container_name=api-go

all:
	@make deps
	@make build
	@make debug

deps:
	@echo "Updating dependencies"
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated successfully"

build:
	@echo "Building the application"
	@echo "Current directory: $(CURDIR)"
	@echo $(wildcard ./internal/config/*)
	go build -tags 'no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' -o ./bin/blog-api -race ./cmd/blog-api/
	copy .\internal\config\config.yaml .\bin
	@echo "Build successful"

run:
	@echo "Running the application"
	./bin/blog-api --debug=false

debug:
	@echo "Debugging the application"
	.\bin\blog-api --debug=true

image:
	@echo "building docker image"
	docker build -t $(image_name) .
	docker run -d -p 8000:8000 --name $(container_name) $(image_name)

run-container:
	@echo "Running the container"
	docker start $(container_name)

stop-container:
	@echo "Stopping the container"
	docker stop $(container_name)

compose:
	@echo "Composing the images"
	docker compose up -d

compose-down:
	@echo "Removing the composed images"
	docker compose down

db-migrate-up:
	goose -dir ./cmd/migrate/migrations postgres "host=localhost port=5432 database=blogs user=postgres password=admin" up

db-migrate-down:
	goose -dir ./cmd/migrate/migrations postgres "host=localhost port=5432 database=blogs user=postgres password=admin" down

swag-gen:
	swag init -g cmd/blog-api/main.go -o docs

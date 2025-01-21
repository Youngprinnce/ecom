build:
	@go build -o bin/ecom cmd/main.go

run: build
	@./bin/ecom

test:
	@go test -v ./...

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

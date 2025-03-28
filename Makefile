.PHONY: run tests tests-json failed-tests create_test_project migrate short-tests short-tests-details generate create_test_project tests-coverage

run:
	@ulimit -n 4096
	@air -c ./config/air.toml -- -lsptesting
build:
	@go build -o bin/timesheets cmd/main.go
tests: generate
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on 

tests-coverage: generate
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -coverprofile=coverage.out
	
short-tests-details: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -short
short-tests: 
	@echo "Running tests..."
	@go test ./... -race -shuffle=on -short
tests-json: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -json 
failed-tests: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -json | jq '.|select(.Action=="fail" and .Test!=null)'
generate:
	@sqlc generate -f ./config/sqlc.yaml
testDb:
	@sqlc generate -f ./config/sqlc.yaml
	@rm -rf ./temp/timesheets.db*
	@migrate -database "sqlite3://./temp/timesheets.db?journal_mode=WAL&foreign_keys=true&cache_size=2000" -path ./internal/db/schema/migrations up
migrate:
	@migrate -database "sqlite3://./timesheets.db?journal_mode=WAL&foreign_keys=true&cache_size=2000" -path db/migrations up
	@sqlc generate -f ./config/sqlc.yaml
create_test_project:
	bash ./scripts/create_test_project.sh


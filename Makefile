.PHONY: tests coverage coverage-html coverage-total linter clean migrate seed tidy run

# This command will run the tests for the project
tests:
	go test -v ./...

# This command will generate a coverage report for the project
coverage:
	go test -cover -coverprofile=coverage.out ./...

# This command will display the coverage report in an HTML file
coverage-html: coverage
	go tool cover -html="coverage.out"
	rm -f coverage.out

# This command will display the coverage report in the terminal
coverage-total: coverage
	go tool cover -func="coverage.out"

# This command will run the linter for the project
linter:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# Remove test files
clean:
	rm -f coverage.out

# Creates database 
migrate:
	go run cmd/migrate/main.go

# Populates the database tables
seed:
	go run cmd/seed/main.go

# Install dependencies
tidy:
	go mod tidy

# Run the api
run:
	go run cmd/main.go

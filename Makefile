# This command will run the tests for the project
.PHONY: tests
tests:
	go test -v ./...

# This command will generate a coverage report for the project
.PHONY: coverage
coverage:
	go test -cover -coverprofile=coverage.out ./...

# This command will display the coverage report in an HTML file
.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html="coverage.out"
	rm -f coverage.out

# This command will display the coverage report in the terminal
.PHONY: coverage-total
coverage-total: coverage
	go tool cover -func="coverage.out"

# This command will run the linter for the project
.PHONY: linter
linter:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# Remove test files
.PHONY: clean
	rm -f coverage.out

# Creates database 
.PHONY: migrate
	go run cmd/migrate/main.go

# Populates the database tables
.PHONY: seed
	go run cmd/seed/main.go

# Install dependencies
.PHONY: tidy
	go mod tidy

# Run the api
.PHONY: run
	go run cmd/main.go

.PHONY: test coverage coverage-html coverage-total linter clean migrate seed tidy run

PKGS := $(shell go list ./... | grep -vE '/test')

# This command will run the test for the project
test:
	go test -v $(PKGS)

# This command will generate a coverage report for the project
coverage:
	go test -cover -coverprofile=coverage.out $(PKGS)

# This command will display the coverage report in an HTML file
coverage-html: coverage
	go tool cover -html="coverage.out"

# This command will display the coverage report in the terminal
coverage-total: coverage
	go tool cover -func=coverage.out | grep total | awk '{print "total coverage: " $$3}'

# This command will run the linter for the project
linter:
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

coverage-section-html:
	@echo "--- Generating full coverage report for section packages... ---"
	@go test -cover -coverprofile=coverage.out ./internal/handler ./internal/service ./internal/repository
	@echo "--- Filtering report for section.go files... ---"
	@(head -n 1 coverage.out && grep "section.go" coverage.out) > coverage-section.out
	@echo "--- Opening filtered HTML report in browser... ---"
	@go tool cover -html="coverage-section.out"
	@rm -f coverage-section.out
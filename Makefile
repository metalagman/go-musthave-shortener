lint:
	@echo "Running linter checks"
	golangci-lint run

test:
	@echo "Running UNIT tests"
	@go clean -testcache
	go test -cover -race -short ./... | { grep -v 'no test files'; true; }

cover-html:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -html=coverage.out

cover:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -func=coverage.out

generate:
	@echo "Generating mocks"
	go generate ./...

.PHONY: build
build: build-shortener

build-shortener:
	@echo "Building the shortener app to the bin dir"
	CGO_ENABLED=1 go build -o ./bin/shortener ./cmd/shortener/*.go

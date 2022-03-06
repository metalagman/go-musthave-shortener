lint:
	@echo "Running linter checks"
	golangci-lint run

test:
	@echo "Running UNIT tests"
	@go clean -testcache
	go test -cover -race -short ./... | { grep -v 'no test files'; true; }

.PHONY: build
build: build-shortener

build-shortener:
	@echo "Building the shortener app to the bin dir"
	CGO_ENABLED=1 go build -o ./bin/shortener ./cmd/shortener/*.go

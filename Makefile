all: build

.PHONY: build check lint test run

check: test

lint:
	gofmt -w .
	go vet ./...

test: lint
	go test -v -coverprofile=coverage.out .

run: check
	go run ./cmd/ssd1305

coverage.html: coverage.out
	go tool cover -html=$< -o $@

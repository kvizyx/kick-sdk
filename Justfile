# Reference: https://just.systems/man/en

coverage_dir := ".coverage"

lint:
	golangci-lint run --timeout 5m --config .golangci.yaml

install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6

test:
	mkdir -p {{ coverage_dir }}
	go test -coverprofile={{ coverage_dir }}/coverage.out ./...
	go tool cover -html={{ coverage_dir }}/coverage.out -o ./{{ coverage_dir }}/coverage.html

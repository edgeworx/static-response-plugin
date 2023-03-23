default: lint test

lint:
	golangci-lint run -v

test:
	go test -v -cover ./...
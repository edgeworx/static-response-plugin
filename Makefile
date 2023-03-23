default: lint test

lint:
	go install github.com/mgechev/revive@latest
	revive ./...

test:
	go test -v -cover ./...
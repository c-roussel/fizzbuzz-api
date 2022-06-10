lint:
	go vet ./...
	gofmt -s -w .
	go mod tidy

test: lint
	go test ./...

swag:
	swag init -g cmd/server/main.go --output docs/swagger/

build: lint swag
	go build -ldflags="-s -w" ./cmd/server

clean:
	go clean
	if [ -f "./server" ]; then \
		rm ./server; \
	fi

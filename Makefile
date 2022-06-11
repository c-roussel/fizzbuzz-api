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

build_docker_image: build
	GIT_HASH=`git rev-parse --verify HEAD`
	docker build \
		--build-arg GIT_HASH=$GIT_HASH \
		--no-cache \
		-t fizzbuzz-api \
		.

docker: build_docker_image clean

clean:
	go clean
	if [ -f "./server" ]; then \
		rm ./server; \
	fi

# Context

`fizzbuzz-api` is an HTTP server allowing users to try every possible fizzbuzz variants.

The original `Fizz Buzz` algorithm consists in writing all numbers from 1 to 100, and just
 replacing all multiples of 3 by `fizz`, all multiples of 5 by `buzz`, and all multiples of 15 by `fizzbuzz`.

This HTTP server serves a `GET /fizzbuzz`route responding with any alternative of this algorithm.

# Requirements

- [echo-swagger](https://github.com/swaggo/echo-swagger)

# Installation

Several ways to run this API:

- Build the binary using `make build`. Then run the generated `server` binary.

- `go run cmd/server/main.go`

- Build the `fizzbuzz-api` docker image using `make docker`.
Then run `docker run -p $YOUR_PORT:3000 fizzbuzz-api`.

# Routes

Feel free to read `docs/swagger/swagger.yml`.

You also can run the server and reach the `/swagger/index.html` endpoint.

# Contributing

Make sure: 
- Code complies with `make lint` and `make test`.
- Tests are done using [go-testdeep](https://github.com/maxatome/go-testdeep).
- Public function & structures are documented.
- HTTP handler functions are documented using [echo-swagger](https://github.com/swaggo/swag#declarative-comments-format).

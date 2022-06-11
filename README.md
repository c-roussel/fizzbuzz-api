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

# Configuration

Envrionment variables:

- `FIZZBUZZ_MAX_LIMIT`: integer that will limit the maximum `limit` on /fizzbuzz route.

# Monitoring

This API serves a prometheus endpoint on `GET /mon/metrics`.

You may install [prometheus](https://prometheus.io/download/) and run it:

```
prometheus --config.file=prometheus.yml
```

Note that running the prometheus server only makes sense if the API is also running on the same host.

Otherwise, you would need to update the targets in `prometheus.yml`.

# Contributing

Make sure: 
- Code complies with `make lint` and `make test`.
- Tests are done using [go-testdeep](https://github.com/maxatome/go-testdeep).
- Public function & structures are documented.
- HTTP handler functions are documented using [echo-swagger](https://github.com/swaggo/swag#declarative-comments-format).

[![GoDoc Reference](https://pkg.go.dev/badge/github.com/sf9v/kalupi)](https://pkg.go.dev/github.com/sf9v/kalupi)
[![Go Report Card](https://goreportcard.com/badge/github.com/sf9v/kalupi)](https://goreportcard.com/report/github.com/sf9v/kalupi)

# Kalupi

Kalupi, a wallet service built with [go-kit](https://github.com/go-kit/kit).

## Features

- [Double-entry accounting](https://en.wikipedia.org/wiki/Double-entry_bookkeeping)
- Modular and extensible design
- Built with [go-kit](https://github.com/go-kit/kit)!

## Limitations

- Only payments within the same currency is supported
- Only USD is supported, other currencies can be supported with relative ease

## Documentation

The REST API documentation is located at [docs/api.md](/docs/api.md).

## Build

Build the server:

```sh
$ go build -v -ldflags "-w -s" -o ./cmd/kalupi ./cmd/kalupi
```

Run the server:

```sh
$ DSN=<postgres connection string> ./cmd/kalupi
```

## Docker

The container image is hosted on [docker hub](https://hub.docker.com/repository/docker/stevenferrer/kalupi).

Run using docker:
```sh
$ docker run -p 8000:8000 \
	-e DSN=<postgres connection string> \
	stevenferrer/kalupi:0.1.0-rc1
```

## Development

Requirements:
- [Go](https://golang.org/)
- [Postgres](http://postgresql.org/)

Clone the repository:

```sh
$ git clone git@github.com:sf9v/kalupi.git
```

Setup test database:

```sh
$ docker run --name kalupi-test-db \
	-d --rm -p 5432:5432 \
	-e POSTGRES_PASSWORD=postgres \
	postgres:12
```

Run the tests:

```sh
$ go test -v -cover -race ./...
```

## Shoulders of the giants

The double-entry accounting implementation in this project is heavily based on the [ideas](https://stackoverflow.com/questions/59432964/relational-data-model-for-double-entry-accounting) of [PerformanceDBA](https://stackoverflow.com/users/484814/performancedba) and deserves most of the credit.

## Contributing

All contributions are welcome! Please feel free to [open an issue](https://github.com/sf9v/kalupi/issues/new) or [make a pull request](https://github.com/sf9v/kalupi/pulls).

## License

[MIT](LICENSE)
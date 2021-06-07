GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOOS ?=linux"
GOARCH ?=amd64

.PHONY: postgres
postgres:
	docker rm -f postgres || true
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -d --rm -p 5432:5432 docker.io/library/postgres:12
	docker exec -it postgres bash -c 'while ! pg_isready; do sleep 1; done;'

.PHONY: test
test:
	go test -v -cover -race ./...
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOOS ?=linux"
GOARCH ?=amd64

.PHONY: build
build:
	go build -v -ldflags "-w -s" -o ./cmd/kalupi ./cmd/kalupi

.PHONY: test
test:
	go test -v -cover -race ./...


.PHONY: db
db: test-db dev-db

.PHONY: test-db
test-db:
	docker rm -f kalupi-test-db || true
	docker run --name kalupi-test-db \
		-d --rm -p 5432:5432 \
		-e POSTGRES_PASSWORD=postgres \
		postgres:12
	docker exec -it kalupi-test-db bash -c 'while ! pg_isready; do sleep 1; done;'

.PHONY: dev-db
dev-db:
	docker rm -f kalupi-dev-db || true
	docker run --name kalupi-dev-db \
		-d --rm -p 5433:5432 \
		-e POSTGRES_USER=kalupi \
		-e POSTGRES_PASSWORD=kalupi \
		-e POSTGRES_DB=kalupi \
		postgres:12
	docker exec -it kalupi-dev-db bash -c 'while ! pg_isready; do sleep 1; done;'


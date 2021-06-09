FROM golang:1.16-alpine as build

RUN apk update && apk upgrade
RUN apk add --no-cache wget

WORKDIR /src
COPY . /src

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /build/kalupi ./cmd/kalupi

FROM scratch

COPY --from=build /build/kalupi /bin/

ENTRYPOINT ["/bin/kalupi"]
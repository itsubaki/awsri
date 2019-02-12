SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

.PHONY: test
test:
	mkdir -p /var/tmp/hermes/{pricing,costexp}
	GOCACHE=off go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize* -timeout 20m
	GOCACHE=off go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

install:
	-rm ${GOPATH}/bin/hermes
	go install -ldflags "${LDFLAGS}" github.com/itsubaki/hermes

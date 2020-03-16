SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

install:
	-rm ${GOPATH}/bin/hermes
	GO111MODULE=on go mod tidy
	GO111MODULE=on go install -ldflags "${LDFLAGS}" github.com/itsubaki/hermes

.PHONY: test
test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

runmysql:
	set -x
	-docker pull mysql
	-docker stop mysql
	-docker rm mysql
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
	# mysql -h127.0.0.1 -P3306 -uroot -psecret

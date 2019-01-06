SHELL := /bin/bash

serialize:
	mkdir -p ${GOPATH}/src/github.com/itsubaki/awsri/internal/_serialized/awsprice
	mkdir -p ${GOPATH}/src/github.com/itsubaki/awsri/internal/_serialized/costexp
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize* -timeout 20m

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

SHELL := /bin/bash

.PHONY: test
test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize*
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

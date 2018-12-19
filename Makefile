SHELL := /bin/bash


# make test AWS_ACCOUNT_ID=123456789012
.PHONY: test
test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize*
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

SHELL := /bin/bash

serialize:
	mkdir -p /var/tmp/hermes/awsprice
	mkdir -p /var/tmp/hermes/costexp
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize* -timeout 20m

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

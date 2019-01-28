SHELL := /bin/bash

test:
	mkdir -p /var/tmp/hermes/{awsprice,costexp,reserved}
	GOCACHE=off go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v -run TestSerialize* -timeout 20m
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

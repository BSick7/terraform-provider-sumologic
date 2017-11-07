SHELL := /bin/bash

.PHONY: deps vet test

deps:
	go get -u github.com/Masterminds/glide
	go get github.com/jstemmer/go-junit-report
	glide install

vet:
	glide nv | xargs go vet

test:
	set -o pipefail;glide nv \
		| xargs go test -v \
		| tee /dev/tty \
		| go-junit-report > unit-tests.xml

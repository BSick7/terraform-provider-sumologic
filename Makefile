SHELL := /bin/bash

.PHONY: vet build test

deps:
	go get -u github.com/Masterminds/glide
	glide install

vet:
	go vet ./...

build:
	go build

install:
	go install

test:
	go test -v ./... \
		| tee /dev/tty \
		| go-junit-report > unit-tests.xml

release:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr
	gox -os "linux darwin windows" -arch "amd64 386" -ldflags "-X main.Version=`cat VERSION`" -output="dist/terraform-provider-sumologic{{.OS}}_{{.Arch}}"
	ghr -t $$GITHUB_TOKEN -u BSick7 -r terraform-provider-sumologic --replace `cat VERSION` dist/

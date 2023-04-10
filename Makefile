MODULE_NAME := myapp
SRC := $(shell find . -name '*.go')
BUILD_DIR := ./build
BINARY_NAME := $(BUILD_DIR)/$(MODULE_NAME) 

.PHONY: all clean build test coverage vet lint docker-build docker-run docker-push

all: build

build:
	go build -v -ldflags="-X 'version.Version=v1.0.0' -X 'version.GitCommit=$(shell git rev-parse --short=8 HEAD)' -X 'build.User=$(shell id -u -n)' -X 'build.Time=$(shell date)'" -o $(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

unit-test:
	@echo "Running unit tests"
	go test -v $(shell go list ./... | grep -v /test) \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./unit-test-coverage.out

api-test:
	@echo "Running api tests"
	go test -v ./test/ \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./api-test-coverage.out

cov-html:
	go tool cover -html=api-test-coverage.out -html=unit-test-coverage.out -o merged-coverage.html

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	go vet $(SRC)

lint:
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...

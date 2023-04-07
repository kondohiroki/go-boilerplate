MODULE_NAME := myapp
SRC := $(shell find . -name '*.go')
BUILD_DIR := ./bin
BINARY_NAME := $(BUILD_DIR)/$(MODULE_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: all clean build test coverage vet lint docker-build docker-run docker-push

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(SRC)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

unit-test:
	@echo "Running unit tests"
	@$(GOTEST) -v ./... \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./unit-test-coverage.out

api-test:
	@echo "Running api tests"
	@$(GOTEST) -v ./test/... \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./api-test-coverage.out

cov-html:
	go tool cover -html=api-test-coverage.out -html=unit-test-coverage.out -o merged-coverage.html

coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	$(GOCMD) vet $(SRC)

lint:
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...

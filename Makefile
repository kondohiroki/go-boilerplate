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

test:
	$(GOTEST) -v ./... \
	-count=1 \
	-cover \
	-coverprofile=./unit-test-coverage.out

coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	$(GOCMD) vet $(SRC)

lint:
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...

docker-build:
	docker build -t $(MODULE_NAME) .

docker-run:
	docker run -p 8080:8080 $(MODULE_NAME)

docker-push:
	docker tag $(MODULE_NAME) registry.example.com/$(MODULE_NAME)
	docker push registry.example.com/$(MODULE_NAME)

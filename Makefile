MODULE_NAME := myapp
SRC := $(shell find . -name '*.go')
BUILD_DIR := ./build
BINARY_NAME := $(BUILD_DIR)/$(MODULE_NAME)
SONAR_HOST_URL := https://sonarcloud.io
SONAR_SECRET := $(shell cat .sonar.secret)
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
CHANGE_TARGET := $(shell git rev-parse --abbrev-ref --symbolic-full-name @{u} | sed 's/.*\///')
# CHANGE_ID := $(shell git rev-parse --short=8 HEAD)
CHANGE_ID := $(shell whoami)

.PHONY: all

build:
	go build -v -ldflags="-X 'version.Version=v1.0.0' -X 'version.GitCommit=$(shell git rev-parse --short=8 HEAD)' -X 'build.User=$(shell id -u -n)' -X 'build.Time=$(shell date)'" -o $(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

unit-test:
	@echo "Running unit tests"
	rm -f unit-test-coverage.out && \
	go test -v $(shell go list ./... | grep -v /test) \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./unit-test-coverage.out

api-test:
	@echo "Running api tests"
	rm -f api-test-coverage.out && \
	go test -v ./test/ \
	-count=1 \
	-cover \
	-coverpkg=./... \
	-coverprofile=./api-test-coverage.out

unit-test-xml:
	@echo "Running unit tests"
	rm -f unit-test-report.xml && \
	go test -v 2>&1 $(shell go list ./... | grep -v /test) \
	-count=1 \
	| go-junit-report -set-exit-code > unit-test-report.xml

api-test-xml:
	@echo "Running api tests"
	rm -f api-test-report.xml && \
	go test -v 2>&1 ./test/ \
	-count=1 \
	| go-junit-report -set-exit-code > api-test-report.xml

cov-html:
	go tool cover -html=api-test-coverage.out -html=unit-test-coverage.out -o merged-coverage.html

sonarqube-pr:
	rm -rf .scannerwork && \
	sonar-scanner \
		-Dsonar.host.url="$(SONAR_HOST_URL)" \
		-Dsonar.working.directory=".scannerwork" \
		-Dsonar.pullrequest.key="$(CHANGE_ID)" \
		-Dsonar.pullrequest.branch="$(BRANCH_NAME)" \
		-Dsonar.pullrequest.base="$(CHANGE_TARGET)" \
		-Dsonar.login="$(SONAR_SECRET)"

sonarqube-branch:
	rm -rf .scannerwork && \
	sonar-scanner \
		-Dsonar.host.url="$(SONAR_HOST_URL)" \
		-Dsonar.working.directory=".scannerwork" \
		-Dsonar.branch.name="$(BRANCH_NAME)" \
		-Dsonar.login="$(SONAR_SECRET)"

coverage:
	make unit-test && make api-test

vet:
	go vet $(SRC)

lint:
	go get golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...

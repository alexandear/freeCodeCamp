MAKEFILE_PATH := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PATH := $(MAKEFILE_PATH):$(PATH)

export GOBIN := $(MAKEFILE_PATH)/bin

all: clean generate build test

clean:
	@echo clean
	@go clean

build:
	@echo build
	@go build -o $(GOBIN)/messageboard

TEST_PKGS = $(shell go list ./... | grep -v /test)

.PHONY: test
test:
	@echo test
	@go test -race -count=1 -v $(TEST_PKGS)

generate:
	@echo generate
	@go generate ./...

IMAGE = news-api

install-tools:
	@echo install tools from tools.go
	@cat tools.go | grep _ | awk -F '"' '{print $$2}' | xargs -tI % go install %

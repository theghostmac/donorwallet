# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run

# Project parameters
BINARY_NAME=donorwallet
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_FILE=./cmd/main.go

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_FILE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GORUN) $(MAIN_FILE)

deps:
	$(GOCMD) mod download

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_FILE)

# Convenience commands
start: build
	./$(BINARY_NAME) &

stop:
	pkill -f $(BINARY_NAME)

.PHONY: all test clean run deps build-linux start stop

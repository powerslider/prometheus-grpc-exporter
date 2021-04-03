BUF_VERSION:=0.41.0
GOLANGCI_VERSION:=1.39.0
PROJECT_NAME:=prometheus-grpc-exporter
GOPATH_BIN:=$(shell go env GOPATH)/bin

all: clean lint generate build-server build-client

install:
	# Install protobuf compilation plugins.
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

	# Install buf tool for protobuf stub generation, linting, etc.
	curl -sSfL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "${GOPATH_BIN}/buf" && chmod +x "${GOPATH_BIN}/buf"

	# Install golangci-lint for go code linting.
	curl -sSfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | \
		sh -s -- -b ${GOPATH_BIN} v${GOLANGCI_VERSION}


proto-lint:
	buf lint

generate:
	buf generate

lint:
	golangci-lint run --config=.golangci.yml

build-server:
	@echo ">>> Building ${PROJECT_NAME} gRPC server..."
	go build -o bin/server cmd/server/main.go

build-client:
	@echo ">>> Building ${PROJECT_NAME} gRPC client..."
	go build -o bin/client cmd/client/main.go

clean:
	@echo ">>> Removing binaries..."
	@rm -rf bin/*

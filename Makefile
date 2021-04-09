BUF_VERSION:=0.41.0
GOLANGCI_VERSION:=1.39.0
PROJECT_NAME:=prometheus-grpc-exporter
GOPATH_BIN:=$(shell go env GOPATH)/bin

all: clean lint proto-lint build-server build-client build-remote-receiver

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
	@echo ">>> Performing golang code linting.."
	golangci-lint run --config=.golangci.yml

build-server:
	@echo ">>> Building ${PROJECT_NAME} API server..."
	go build -o bin/server cmd/prometheus-api-server/main.go

build-client:
	@echo ">>> Building ${PROJECT_NAME} gRPC client..."
	go build -o bin/client cmd/prometheus-grpc-client/main.go

build-remote-receiver:
	@echo ">>> Building ${PROJECT_NAME} prometheus-remote-receiver..."
	go build -o bin/remote-receiver cmd/prometheus-remote-receiver/main.go

clean:
	@echo ">>> Removing old binaries..."
	@rm -rf bin/*

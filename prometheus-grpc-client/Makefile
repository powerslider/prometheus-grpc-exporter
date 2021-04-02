GRPC_GW_GOOGLEAPIS=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GOGO_PROTOBUF=${GOPATH}/src/github.com/gogo/protobuf

all: clean test build

build:
	@echo ">>> Building Application..."
	go build -o bin/prometheus-grpc-client cmd/prometheus-grpc-client/main.go

test:
	@echo ">>> Running Unit Tests..."
	go test -race ./...

cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	@echo ">>> Removing binaries..."
	@rm -rf bin/*

generate:
	protoc \
		-I proto \
		-I ${GRPC_GW_GOOGLEAPIS} \
		-I ${GOGO_PROTOBUF} \
		--proto_path=${HOME}/go/src/github.com/powerslider/prometheus-grpc-client \
		--gogo_out=plugins=grpc,paths=source_relative:./proto \
		--grpc-gateway_out=allow_patch_feature=false,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,logtostderr=true:./proto \
		proto/*.proto

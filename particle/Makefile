.PHONY: all protos
default: all

export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOBIN?=$(BIN)
export GO=$(shell which go)
export BUILD=cd $(ROOT) && $(GO) install -v -ldflags "-s"

all:
	$(BUILD) ./cmd/...

protos:
	protoc -I pkg/gateway/ pkg/gateway/kafka.proto --go_out=pkg/gateway


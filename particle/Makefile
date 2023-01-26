.PHONY: all protos
default: all

all:
	go install -v -ldflags "-s" ./cmd/...

protos:
	protoc -I pkg/gateway/ pkg/gateway/kafka.proto --go_out=pkg/gateway


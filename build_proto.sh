#!/bin/bash
cd ./proto/
# export GO111MODULE=on  # Enable module mode
# go get github.com/golang/protobuf/protoc-gen-go@v1.3
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I . rpc.proto --go_out=plugins=grpc:.
protoc -I . cluster.proto --go_out=plugins=grpc:.
echo "Done"
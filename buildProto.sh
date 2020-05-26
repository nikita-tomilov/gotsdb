#!/bin/bash
cd ./proto/
# export GO111MODULE=on  # Enable module mode
# go get github.com/golang/protobuf/protoc-gen-go@v1.3
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I . testInteraction.proto --go_out=plugins=grpc:.
echo "Done"
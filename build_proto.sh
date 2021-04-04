#!/bin/bash
export GO111MODULE=on  # Enable module mode
# go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
# go get -u google.golang.org/grpc
# sudo apt install protobuf-c*
export PATH="$PATH:$(go env GOPATH)/bin"

SRC_REL_DIR="."
PROTO_REL_DIR=$SRC_REL_DIR/proto

for protofile in rpc.proto cluster.proto
do
  protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $PROTO_REL_DIR/$protofile
done

echo "Done"
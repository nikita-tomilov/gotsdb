#!/bin/bash
set -x
echo Launching tests
go test -v ./... && echo "tests OK" || echo "test FAILURE"

#!/bin/bash
. VERSION
docker build . -t gotsdb-srv:$VERSION
# remove unneeded ones
docker image rm $(docker image ls --filter dangling=true -q)
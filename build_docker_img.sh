#!/bin/bash
docker build . -t gotsdb-srv:v0.2
# remove unneeded ones
docker image rm $(docker image ls --filter dangling=true -q)
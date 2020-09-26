#!/bin/bash
docker run -v `pwd`/config:/config -p 5300:5300 -p 5123:5123 gotsdb-srv:latest
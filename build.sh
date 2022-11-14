#!/usr/bin/env bash
docker ps -a | egrep 'user-management-api' | awk '{print $1}'| xargs docker rm -f
docker build -t user-management-api .
docker run  -itd  --name user-management-api -p 80:80 user-management-api
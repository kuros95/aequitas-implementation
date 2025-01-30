#!/bin/bash

container_id=$(echo$(docker ps -a | grep client | cut -d ' '  -f 1))

docker cp $container_id:/builder/client.log ./client.log
docker cp $container_id:/builder/dump.log ./dump.log
#!/bin/bash

rm -rf edge*
rm -rf vertex*
rm -rf *.json
rm -rf *.log
docker stop $(docker ps -aq)
docker container prune -f
docker rmi -f vertex1 edge1 edge2 vertex3 
docker network rm graph 

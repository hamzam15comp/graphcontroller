#!/bin/bash

rm -rf edge*
rm -rf vertex*
rm -rf *.json
docker stop $(docker ps -aq)
docker container prune -f
docker rmi -f vertex1 vertex2 edge1 edge2 vertex3 vertex4
docker network rm graph 

#!/bin/bash

go run creategraph.go

sleep 120 

docker cp vertex1:/go/apperrors.log /hamzam/imglog/v1.log
docker cp vertex2:/go/apperrors.log /hamzam/imglog/v2.log
docker cp vertex3:/go/apperrors.log /hamzam/imglog/v3.log

go run /hamzam/dockercli/tests/replaceV2V4.go

sleep 30

docker cp vertex1:/go/apperrors.log /hamzam/imglog/v1.log
docker cp vertex4:/go/apperrors.log /hamzam/imglog/v4.log
docker cp vertex3:/go/apperrors.log /hamzam/imglog/v3.log

#!/bin/bash

docker build --build-arg COMMIT=`git rev-parse HEAD` -t mlesniak/go-demo .

# See https://docs.travis-ci.com/user/docker/#pushing-a-docker-image-to-a-registry
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push mlesniak/go-demo

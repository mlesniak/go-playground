#!/bin/bash

# Build new container.
docker build --build-arg COMMIT=`git rev-parse HEAD` -t mlesniak/go-demo .
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push mlesniak/go-demo


# Deploy to kubernetes by setting new image.
sed -i -e 's|KUBE_CA_CERT|'"${KUBE_CA_CERT}"'|g'       deployments/kube-config.yaml
sed -i -e 's|KUBE_ID|'"${KUBE_ID}"'|g'                 deployments/kube-config.yaml
sed -i -e 's|KUBE_TOKEN|'"${KUBE_TOKEN}"'|g'           deployments/kube-config.yaml

# We currently do not use $TRAVIS_BUILD_ID intentionally.
docker run --rm \
    -v $(pwd)/deployments/kube-config.yaml:/.kube/config\
    bitnami/kubectl:latest \
    set image deployment/go-demo go-demo=${DOCKER_USERNAME}/go-demo:latest


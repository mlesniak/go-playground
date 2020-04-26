[![Build Status](https://travis-ci.com/mlesniak/go-playground.svg?branch=master)](https://travis-ci.com/mlesniak/go-playground)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mlesniak_go-playground&metric=alert_status)](https://sonarcloud.io/dashboard?id=mlesniak_go-playground)
[![Code of Conduct](https://img.shields.io/badge/%E2%9D%A4-code%20of%20conduct-orange.svg?style=flat)](CODE_OF_CONDUCT.md)

# Overview

This is a playground for misc. go frameworks usable in a production system.

    http mlesniak.dev/api numbers:=10


## Add secret logging token

    echo -n "TOKEN"|kubectl create secret generic sematext-token --from-file=token=/dev/stdin


## Enable filebeat kubernetes authentication

Given the error message

    Failed to list *v1.Pod: pods is forbidden: User "system:serviceaccount:default:default" cannot list resource "pods" in API group "" at the cluster scope

execute

    kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default


# Next steps

- [ ] Add filebeat for log submission
- [ ] Use configmap for filebeat.yml (see [https://github.com/elastic/beats/blob/master/deploy/kubernetes/filebeat-kubernetes.yaml](here))
- [ ] Add persistence storage for filebeat over PVC/BlockStorage in K8s
- [ ] add nginx for subdomains in k8
- [ ] JWT middleware in echo?
- [ ] Think about keycloak?
- [ ] Add integration tests
- [ ] Extract log package
- [ ] Extract computation service
- [ ] automatic K8 deployment over Travis CI

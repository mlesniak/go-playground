[![Build Status](https://travis-ci.com/mlesniak/go-playground.svg?branch=master)](https://travis-ci.com/mlesniak/go-playground)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mlesniak_go-playground&metric=alert_status)](https://sonarcloud.io/dashboard?id=mlesniak_go-playground)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlesniak/go-playground)](https://goreportcard.com/report/github.com/mlesniak/go-playground)
[![Code of Conduct](https://img.shields.io/badge/%E2%9D%A4-code%20of%20conduct-orange.svg?style=flat)](CODE_OF_CONDUCT.md)

# Overview

This is a playground for misc. go frameworks usable in a production system.

    http https://api.mlesniak.dev/api number:=30


## Add loadbalancer forwarding rule in Digital Ocean for HTTPS access

![screenshot](docs/loadbalancer-rules.png)

## Add secrets

Secret handling is explicitliy manual

    echo -n "TOKEN"|kubectl create secret generic sematext-token --from-file=token=/dev/stdin
    echo -n "PASSWORD"|kubectl create secret generic keycloak-passowrd --from-file=password=/dev/stdin


## Enable filebeat kubernetes authentication

Given the error message

    Failed to list *v1.Pod: pods is forbidden: User "system:serviceaccount:default:default" cannot list resource "pods" in API group "" at the cluster scope

execute

    kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

Note: the creation of the role is now executed in the `yaml` file, but since it was quite difficult to find information about this problem
online I leave it here for now.

## Example authentication with keycloak

We use our app as a proxy

    # docker run --name keycloak -p 8081:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:9.0.3
    kubectl port-forward service/keycloak-service 8081:8080

    export A=$(http POST :8080/api/login username=demo password=demo)
    export T=$(echo $A|jq -r .accessToken)
    export R=$(echo $A|jq -r .refreshToken)
    http POST :8080/api number:=10 Authorization:"Bearer $T"
    http -v POST :8080/api/logout Authorization:"Bearer $T" refreshToken=$R username=demo password=demo


# Next steps

- [ ] Enable environment support and reconfigure application deployment for keycloak
- [ ] Make keycloak available under keycloak.mlesniak.dev
- [ ] Update docker-compose file to use keycloak
- [ ] Add integration tests
- [ ] Add swagger
- [ ] Add mongo database support


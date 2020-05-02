[![Build Status](https://travis-ci.com/mlesniak/go-playground.svg?branch=master)](https://travis-ci.com/mlesniak/go-playground)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mlesniak_go-playground&metric=alert_status)](https://sonarcloud.io/dashboard?id=mlesniak_go-playground)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlesniak/go-playground)](https://goreportcard.com/report/github.com/mlesniak/go-playground)
[![Code of Conduct](https://img.shields.io/badge/%E2%9D%A4-code%20of%20conduct-orange.svg?style=flat)](CODE_OF_CONDUCT.md)

# Overview

This is a playground for misc. go frameworks usable in a production system.

    http https://api.mlesniak.dev/api number:=30


## Add loadbalancer forwarding rule in Digital Ocean for HTTPS access

![screenshot](docs/loadbalancer-rules.png)

## Add secret logging token

    echo -n "TOKEN"|kubectl create secret generic sematext-token --from-file=token=/dev/stdin


## Enable filebeat kubernetes authentication

Given the error message

    Failed to list *v1.Pod: pods is forbidden: User "system:serviceaccount:default:default" cannot list resource "pods" in API group "" at the cluster scope

execute

    kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

Note: the creation of the role is now executed in the `yaml` file, but since it was quite difficult to find information about this problem
online I leave it here for now.

## Example authentication with jwt

    # Beforehand, create K8 secret
    echo -n "..."|kubectl create secret generic jwt-token --from-file=token=/dev/stdin

    export T=$(jwt encode --secret beihaeliech1thooya5Peecu1mi3SheithooFee5cheesheeG1eeGh3Duu9yeiph -P user=1234)
    http :8080/api/version Authorization:"Bearer $T"

# Next steps

- [ ] Think about keycloak integration?
- [ ] Add integration tests


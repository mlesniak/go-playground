[![Build Status](https://travis-ci.com/mlesniak/go-playground.svg?branch=master)](https://travis-ci.com/mlesniak/go-playground)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mlesniak_go-playground&metric=alert_status)](https://sonarcloud.io/dashboard?id=mlesniak_go-playground)
[![Code of Conduct](https://img.shields.io/badge/%E2%9D%A4-code%20of%20conduct-orange.svg?style=flat)](CODE_OF_CONDUCT.md)

# Overview

This is a playground for misc. go frameworks usable in a production system.

    http mlesniak.dev/api numbers:=10


# Add secret logging token

    echo -n "TOKEN"|kubectl create secret generic sematext-token --from-file=token=/dev/stdin

# Next steps

- [ ] Add filebeat for log submission
- [ ] add nginx for subdomains in k8
- [ ] JWT middleware in echo?
- [ ] Think about keycloak?
- [ ] Add integration tests
- [ ] Extract log package
- [ ] Extract computation service
- [ ] automatic K8 deployment over Travis CI

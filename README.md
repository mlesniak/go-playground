[![Build Status](https://travis-ci.com/mlesniak/go-playground.svg?branch=master)](https://travis-ci.com/mlesniak/go-playground)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mlesniak_go-playground&metric=alert_status)](https://sonarcloud.io/dashboard?id=mlesniak_go-playground)
[![Code of Conduct](https://img.shields.io/badge/%E2%9D%A4-code%20of%20conduct-orange.svg?style=flat)](CODE_OF_CONDUCT.md)

# Overview

This is a playground for miscellaneous go frameworks usable in a production system.

# Next steps

- [ ] docker-compose configuration
- [ ] local filebeat docker file
- [ ] Enable log rotation
- [ ] Submit logs to Sematext
- [ ] manual K8s deployment
- [ ] automatic K8 deployment over Travis CI

# Temporary code

    docker run --rm -it \
        -v $(pwd)/deployments/filebeat/filebeat.yml:/usr/share/filebeat/filebeat.yml \
        -v $(pwd)/logs:/logs \
        docker.elastic.co/beats/filebeat:7.6.2

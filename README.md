Sourcehawk FaaS
---------------
 
[![Build Status](https://github.com/optum/sourcehawk-faas/workflows/Maven%20CI/badge.svg)](https://github.com/optum/sourcehawk-faas/actions) 
[![Sourcehawk Scan](https://github.com/optum/sourcehawk-faas/workflows/Sourcehawk%20Scan/badge.svg)](https://github.com/optum/sourcehawk-faas/actions) 
![OSS Lifecycle](https://img.shields.io/osslifecycle/optum/sourcehawk-faas) 

### Docker Images

![Scan Docker Image Version](https://img.shields.io/docker/v/optumopensource/sourcehawk-openfaas-scan) 
![Validate Config Docker Image Version](https://img.shields.io/docker/v/optumopensource/sourcehawk-openfaas-validate-config) 


Sourcehawk Function as a Service offerings.

## Current Implementations

1. [OpenFaaS](openfaas/README.md)

Future implementations may include `Microsoft Azure` or `AWS Lambda`

## Development

### Building
```sh
./mvnw clean install
```
Alidns-Webhook
---

[![Go Report Card](https://goreportcard.com/badge/github.com/wjiec/alidns-webhook)](https://goreportcard.com/report/github.com/wjiec/alidns-webhook)
[![GitHub license](https://img.shields.io/github/license/wjiec/alidns-webhook.svg)](https://github.com/wjiec/alidns-webhook/blob/main/LICENSE)

## Overview

alidns-webhook is a generic ACME solver for [cert-manager](https://github.com/cert-manager/cert-manager).

### Quick start

If you have Helm, you can deploy the alidns-webhook with the following command:
```bash
helm upgrade --install alidns-webhook alidns-webhook \
    --repo https://wjiec.github.io/alidns-webhook \
    --namespace cert-manager --create-namespace \
    --set groupName=acme.yourcompany.com
```
It will install the alidns-webhook in the cert-manager namespace, creating that namespace if it doesn't already exist.


### Supported Versions table

Supported versions for the ingress-nginx project mean that we have completed E2E tests, and they are passing for
the versions listed. Ingress-Nginx versions may work on older versions but the project does not make that guarantee.

| Alidns-Webhook version | k8s supported version  | Helm Chart Version |
|------------------------|------------------------|--------------------|
| **v0.1.0**             | 1.27, 1.26, 1.25, 1.24 | 0.1.*              |


## License

[MIT License](https://github.com/wjiec/alidns-webhook/blob/main/LICENSE)

Alidns-Webhook
---

[![Go Report Card](https://goreportcard.com/badge/github.com/wjiec/alidns-webhook)](https://goreportcard.com/report/github.com/wjiec/alidns-webhook)
[![GitHub license](https://img.shields.io/github/license/wjiec/alidns-webhook.svg)](https://github.com/wjiec/alidns-webhook/blob/main/LICENSE)
[![Kubernetes Compatible](https://github.com/wjiec/alidns-webhook/actions/workflows/k8s-compatible.yml/badge.svg)](https://github.com/wjiec/alidns-webhook/actions/workflows/k8s-compatible.yml)

## Overview

alidns-webhook is a generic ACME solver for [cert-manager](https://github.com/cert-manager/cert-manager).

### Quick start

This tutorial will detail how to configure and install the webhook to your cluster with alidns.


#### Install webhook

Before installing this webhook, make sure you have `cert-manager` installed correctly.
If you haven't installed it yet, you can get the installation instructions from the [cert-manager documentation][1].

If you have Helm, you can deploy the alidns-webhook with the following command:
```bash
helm upgrade --install alidns-webhook alidns-webhook \
    --repo https://wjiec.github.io/alidns-webhook \
    --namespace cert-manager --create-namespace \
    --set groupName=acme.yourcompany.com

# Note: If you installed cert-manager via bitnami charts, you need to add the additional
#   `--set certManager.serviceAccountName=cert-manager-controller`
# parameter to specify the ServiceAccount to use.
```

It will install the alidns-webhook in the cert-manager namespace, creating that namespace if it doesn't already exist.

##### Aliyun registry

If you can't get the image directly through DockerHub, you can use Aliyun's image repository
by adding the following parameter to the installation command:
```plain
--set image.repository=registry.cn-hangzhou.aliyuncs.com/wjiec/alidns-webhook
```


#### Configure a issuer

Create this definition locally and update the email address and groupName to your own. Please see more details in [cert-manager configuration][2].

__Ensure the `groupName` matches the config in the webhook.__

```yaml
#
# example-acme-issuer.yaml
#

apiVersion: v1
kind: Secret
metadata:
  name: alidns-secret
  namespace: cert-manager
stringData:
  access-key-id: "Your Access Key Id"
  access-key-secret: "Your Access Key Secret"
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: example-acme
spec:
  acme:
    # The ACME server URL
    server: https://acme-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration
    email: your@example.com # Change ME
    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
      name: example-acme
    solvers:
      - dns01:
          webhook:
            groupName: acme.yourcompany.com # Change ME
            solverName: alidns
            config:
              region: "cn-hangzhou" # Optional
              accessKeyIdRef:
                name: alidns-secret
                key: access-key-id
              accessKeySecretRef:
                name: alidns-secret
                key: access-key-secret
```

Once edited, apply the custom resource:
```bash
kubectl create --edit -f example-acme-issuer.yaml
```


#### Creating Certificate or deploy a TLS Ingress

We can deploy a certificate directly on Ingress, edit the ingress add the annotations:
```yaml
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: foo-example-com
  annotations:
    cert-manager.io/cluster-issuer: "example-acme"
    # cert-manager.io/issuer: "example-acme"
spec:
  tls:
  - hosts:
    - foo.example.com
    secretName: foo-example-com-tls
  rules:
  - host: foo.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              name: http
```

Or we can create a Certificate resource that is to be honored by an issuer which is to be kept up-to-date.
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: star-example-com
spec:
  secretName: star-example-com-tls
  commonName: "example.com"
  dnsNames:
  - "example.com"
  - "*.example.com"
  issuerRef:
    name: example-acme
    kind: ClusterIssuer
    # kind: Issuer
```
Then we can refer to that secrets(`secretName`) in Ingress.


### Supported Versions table

The following table lists the correspondences between alidns-webhook and k8s versions.

| Alidns-Webhook version | k8s supported version              | Helm Chart Version |
|------------------------|------------------------------------|--------------------|
| **v1.0.&ast;**         | 1.31, 1.30, 1.29, 1.28, 1.27, 1.26 | 1.0.*              |
| **v0.1.0**             | 1.31, 1.30, 1.29, 1.28, 1.27, 1.26 | 0.1.*              |


## License

[MIT License](https://github.com/wjiec/alidns-webhook/blob/main/LICENSE)


[1]: https://cert-manager.io/docs/installation/
[2]: https://cert-manager.io/docs/configuration/

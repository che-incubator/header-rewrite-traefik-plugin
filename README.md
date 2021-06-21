# header-rewrite-proxy

Simple reverse proxy that can modify http headers.

### Motivation

This project was created because [openshift/oauth-proxy](https://github.com/openshift/oauth-proxy)
can't put OpenShift token into `Authorization` header for the upstream application, where we need it
for [kube-rbac-proxy](https://github.com/openshift/kube-rbac-proxy) to work properly. It can only
put it into `X-Forwarded-Access-Token` header and thus we need this component to take the token
from `X-Forwarded-Access-Token` and set it into `Authorization` header.

Note: Upstream [oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) can do what we need, so
we don't need this component on Kubernetes. Some guys has made an effort to port this feature into
OpenShift fork, but it was rejected.

Note: There is ongoing work on Traefik to support
plugins (https://github.com/traefik/traefik/issues/7088). Once this is done, we can implement this
functionality as simple Traefik plugin and dedicated header-rewrite-proxy component won't be needed.

### Configuration

```
Usage of ./main:
  -bind string
    	Listen on. (default ":8080")
  -rules string
    	Path to rules file. (default "rules.yaml")
  -upstream string
    	Full url to upstream application. Mandatory!
```

### Rules
Rewrite rules are configured with yaml file in format
(see https://github.com/che-incubator/header-rewrite-proxy/blob/main/pkg/proxy/conf.go):
```.yaml
rules:
- from: X-Header-Key-From
  to: X-Header-Key-To
  prefix: 'Prefix of header value '
  keep-original: false
  keep-original-target: false
```

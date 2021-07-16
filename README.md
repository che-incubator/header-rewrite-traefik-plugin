# header-rewrite-proxy

Traefik plugin that can modify http headers.

### Motivation

This project was created because [openshift/oauth-proxy](https://github.com/openshift/oauth-proxy)
can't put OpenShift token into `Authorization` header for the upstream application, where we need it
for [kube-rbac-proxy](https://github.com/openshift/kube-rbac-proxy) to work properly. It can only
put it into `X-Forwarded-Access-Token` header and thus we need this plugin to take the token
from `X-Forwarded-Access-Token` and set it into `Authorization` header with `Bearer` prefix.

Note: Upstream [oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) can do what we need, so
we don't need this component on Kubernetes. Some guys has made an effort to port this feature into
OpenShift fork, but it was rejected.

### Configuration

Traefik static configuration for local plugin:
```.yaml
...
experimental:
  localPlugins:
    header-rewrite-proxy:
      moduleName: github.com/che-incubator/header-rewrite-proxy
```

Plugin is then configured as a route middleware
```.yaml
http:
  routers:
    route1:
      middlewares: ["headerRewrite"]
  middlewares:
    headerRewrite:
      plugin:
        header-rewrite-proxy:
          from: X-Forwarded-Access-Token
          to: Authorization
          prefix: 'Bearer '
          keepOriginal: false
          keepOriginalTarget: true
```

In this example, when there is `X-Forwarded-Access-Token` header, it's values will be moved
into `Authorization` header with `Bearer ` as a value prefix. The `X-Forwarded-Access-Token` header
will be removed (`keepOriginal: false`) and if there were any values in `Authorization` header, they
will stay there (`keepOriginalTarget: true`).

# header-rewrite-proxy

Simple reverse proxy that can modify http headers.

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
Rewrite rules are configured with yaml file in format:
```.yaml
rules:
- from: X-Header-Key-From
  to: X-Header-Key-To
  prefix: 'Prefix of header value '
  keep-original: false
```

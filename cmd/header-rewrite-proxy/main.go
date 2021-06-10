package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type ProxyHandler struct {
    p *httputil.ReverseProxy
}

type Conf struct {
    upstream string
    bind string
}

func main() {
    conf, err := parseArgs()
    if err != nil {
        panic(err)
    }

    remote, err := url.Parse(conf.upstream)
    if err != nil {
        panic(err)
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    http.Handle("/", &ProxyHandler{proxy})
    log.Printf("Listening on '%s' ...\n", conf.bind)
    err = http.ListenAndServe(conf.bind, nil)
    if err != nil {
        panic(err)
    }
}
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("X-Forwarded-Access-Token")
    if token != "" {
        log.Println("Writing token from 'X-Forwarded-Access-Token' into 'Authorization' header")
        r.Header.Set("Authorization", "Bearer " + token)
    }
    ph.p.ServeHTTP(w, r)
}

func parseArgs() (*Conf, error) {
    conf := Conf{}
    flag.StringVar(&conf.upstream, "upstream", "", "full url to upstream application")
    flag.StringVar(&conf.bind, "bind", ":8080", "listen on")

    flag.Parse()

    if conf.upstream == "" {
        return nil, fmt.Errorf("upstream parameter must be set")
    }

    return &conf, nil
}

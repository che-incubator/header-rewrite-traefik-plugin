package main

import (
    "flag"
    "fmt"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type ProxyHandler struct {
    p *httputil.ReverseProxy
}

type Conf struct {
    upstream  *url.URL
    bind      string
    rulesConf string
    rules     *Rules
}

type Rules struct {
    Rules []Rule `yaml:"rules"`
}

type Rule struct {
    From         string `yaml:"from"`
    To           string `yaml:"to"`
    Prefix       string `yaml:"prefix"`
    KeepOriginal bool   `yaml:"keep-original"`
}

var conf *Conf

func main() {
    parsedConf, err := parseConfig()
    if err != nil {
        panic(err)
    }
    conf = parsedConf

    proxy := httputil.NewSingleHostReverseProxy(conf.upstream)
    http.Handle("/", &ProxyHandler{proxy})
    log.Printf("Listening on '%s' ...\n", conf.bind)
    err = http.ListenAndServe(conf.bind, nil)
    if err != nil {
        panic(err)
    }
}
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, rule := range conf.rules.Rules {
        headerValue := r.Header.Get(rule.From)
        if headerValue != "" {
            if len(rule.Prefix) > 0 {
                headerValue = rule.Prefix + headerValue
            }
            r.Header.Set(rule.To, headerValue)
            if !rule.KeepOriginal {
                r.Header.Del(rule.From)
            }
        }
    }

    ph.p.ServeHTTP(w, r)
}

func parseConfig() (*Conf, error) {
    conf := Conf{}
    upstream := ""
    flag.StringVar(&upstream, "upstream", "", "Full url to upstream application.")
    flag.StringVar(&conf.bind, "bind", ":8080", "Listen on. Default ':8080'")
    flag.StringVar(&conf.rulesConf, "rules", "rules.yaml", "Path to rules file. Default 'rules.yaml'")

    flag.Parse()

    // verify, parse and set upstream
    if upstream == "" {
        return nil, fmt.Errorf("upstream parameter must be set")
    }
    upstreamUrl, err := url.Parse(upstream)
    if err != nil {
        return nil, err
    }
    conf.upstream = upstreamUrl

    // read rules file
    parsedRules, err := readRulesConf(conf.rulesConf)
    if err != nil {
        return nil, err
    }
    conf.rules = parsedRules
    log.Printf("Using rules: '%+v'", conf.rules)

    return &conf, nil
}

func readRulesConf(filename string) (*Rules, error) {
    buf, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    rules := &Rules{}
    err = yaml.Unmarshal(buf, rules)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %v", filename, err)
    }

    return rules, nil
}

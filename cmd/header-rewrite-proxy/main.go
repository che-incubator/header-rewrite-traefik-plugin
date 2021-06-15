//
// Copyright (c) 2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"header-rewrite-proxy/pkg/proxy"
	"io/ioutil"
	"log"
	"net/url"
)

func main() {
	parsedConf, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := proxy.Serve(parsedConf); err != nil {
		log.Fatal(err)
	}
}

func parseConfig() (*proxy.Conf, error) {
	conf := proxy.Conf{}
	upstream := ""
	flag.StringVar(&upstream, "upstream", "", "Full url to upstream application. Mandatory!")
	flag.StringVar(&conf.Bind, "bind", ":8080", "Listen on.")
	flag.StringVar(&conf.RulesConf, "rules", "rules.yaml", "Path to rules file.")

	flag.Parse()

	// verify, parse and set upstream
	if upstream == "" {
		return nil, fmt.Errorf("upstream parameter must be set")
	}
	upstreamUrl, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}
	conf.Upstream = upstreamUrl

	// read rules file
	parsedRules, err := readRulesConf(conf.RulesConf)
	if err != nil {
		return nil, err
	}
	conf.Rules = parsedRules
	log.Printf("Using rules: '%+v'", conf.Rules)

	return &conf, nil
}

func readRulesConf(filename string) (*proxy.Rules, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	rules := &proxy.Rules{}
	err = yaml.Unmarshal(buf, rules)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return rules, nil
}

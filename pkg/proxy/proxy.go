// Package proxy
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
package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Handler struct {
	proxy *httputil.ReverseProxy
	conf  *Conf
}

func Serve(conf *Conf) error {
	reverseProxy := httputil.NewSingleHostReverseProxy(conf.Upstream)
	http.Handle("/", &Handler{proxy: reverseProxy, conf: conf})
	log.Printf("Listening on '%s' ...\n", conf.Bind)
	if err := http.ListenAndServe(conf.Bind, nil); err != nil {
		return err
	}
	return nil
}

func (proxyHandler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rewriteHeaders(&r.Header, proxyHandler.conf.Rules)
	proxyHandler.proxy.ServeHTTP(w, r)
}

func rewriteHeaders(headers *http.Header, rules *Rules) {
    for _, rule := range rules.Rules {
        headerValue := headers.Get(rule.From)
        if headerValue != "" {
            if len(rule.Prefix) > 0 {
                headerValue = rule.Prefix + headerValue
            }
            headers.Set(rule.To, headerValue)
            if !rule.KeepOriginal {
                headers.Del(rule.From)
            }
        }
    }
}

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

import "net/url"

type Conf struct {
	Upstream  *url.URL
	Bind      string
	RulesConf string
	Rules     *Rules
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

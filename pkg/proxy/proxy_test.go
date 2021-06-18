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
	"net/http"
	"testing"
)

func TestHeaderRewritten(t *testing.T) {
	headers := &http.Header{}
	headers.Add("hello", "world")

	rewriteHeaders(headers, &Rules{Rules: []Rule{
		{
			From:         "hello",
			To:           "goodbye",
			Prefix:       "",
			KeepOriginal: true,
		},
	}})

	if headers.Get("goodbye") != "world" {
		t.Errorf("header 'hello' should be rewritten to new key 'goodbye'")
	}

	if headers.Get("hello") != "world" {
	}
}

func TestKeepOriginal(t *testing.T) {
	headers := &http.Header{}
	headers.Add("hello", "world")

	rewriteHeaders(headers, &Rules{Rules: []Rule{
		{
			From:         "hello",
			To:           "goodbye",
			Prefix:       "",
			KeepOriginal: false,
		},
	}})

	if headers.Get("goodbye") != "world" {
		t.Errorf("header 'hello' should be rewritten to new key 'goodbye'")
	}

	if headers.Get("hello") != "" {
		t.Errorf("original header should be removed")
	}
}

func TestPrefix(t *testing.T) {
	headers := &http.Header{}
	headers.Add("hello", "world")

	rewriteHeaders(headers, &Rules{Rules: []Rule{
		{
			From:         "hello",
			To:           "goodbye",
			Prefix:       "prefix ",
			KeepOriginal: true,
		},
	}})

	if headers.Get("goodbye") != "prefix world" {
		t.Errorf("header 'hello' should be rewritten to new key 'goodbye'")
	}

	if headers.Get("hello") != "world" {
		t.Errorf("original header should not be removed or modified")
	}
}

func TestMultipleValuesUnderSameKey(t *testing.T) {
	headers := &http.Header{}
	headers.Add("hello", "world")
	headers.Add("hello", "there")

	rewriteHeaders(headers, &Rules{Rules: []Rule{
		{
			From:         "hello",
			To:           "goodbye",
			Prefix:       "prefix ",
			KeepOriginal: false,
		},
	}})

	headersResult := headers.Values("goodbye")

	if len(headersResult) != 2 {
		t.Errorf("all header values under same key must be rewritten. '%+v'", headersResult)
	}

	if headersResult[0] != "prefix world" {
		t.Errorf("all header values under same key must be rewritten. '%+v'", headersResult)
	}

	if headersResult[1] != "prefix there" {
		t.Errorf("all header values under same key must be rewritten. '%+v'", headersResult)
	}

	if headers.Get("hello") != "" {
		t.Errorf("original header must not be removed when keepOriginal: 'false'")
	}
}

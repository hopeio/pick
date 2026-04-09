/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"net/url"

	stringsx "github.com/hopeio/gox/strings"
	"github.com/valyala/fasthttp"
)

type ArgsSource fasthttp.Args

func (form *ArgsSource) Get(key string) ([]string, bool) {
	var values []string
	(*fasthttp.Args)(form).VisitAll(func(k, v []byte) {
		if string(k) == key {
			values = append(values, stringsx.FromBytes(v))
		}
	})
	return values, len(values) > 0
}

type CtxSource fasthttp.RequestCtx


func (form *CtxSource) Get(key string) (string, bool) {
	v := (*fasthttp.RequestCtx)(form).UserValue(key).(string)
	return v, v != ""
}

type HeaderSource fasthttp.RequestHeader

func (form *HeaderSource) Get(key string) ([]string, bool) {
	var values []string
	(*fasthttp.RequestHeader)(form).VisitAll(func(k, v []byte) {
		if string(k) == key {
			v, _ := url.QueryUnescape(stringsx.FromBytes(v))
			values = append(values, v)
		}
	})
	return values, len(values) > 0
}

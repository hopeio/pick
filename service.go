/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"net/http"

	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
)

var (
	isRegistered = false
)

type Service[T any] interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []T)
}

func Registered() {
	isRegistered = true
	go func() {
		Openapi(apidoc.Dir, "api")
	}()
	apidoc.ApiDoc(http.DefaultServeMux, apidoc.UriPrefix, apidoc.Dir)
}

func Api(f func()) {
	if !isRegistered {
		f()
	}
}

var prefix string

func HandlerPrefix(p string) {
	prefix = p
}

var DefaultMarshaler = httpx.DefaultMarshal

func Marshaler(marshaler httpx.MarshalFunc) {
	DefaultMarshaler = marshaler
	httpx.DefaultMarshal = marshaler
}

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"reflect"

	"github.com/hopeio/gox/kvstruct"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/valyala/fasthttp"
)

type ArgsSource fasthttp.Args

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *ArgsSource) TrySet(value reflect.Value, field *reflect.StructField, key string, opt *kvstruct.Options) (isSet bool, err error) {
	return kvstruct.SetValueByGetter(value, field, form, key, opt)
}

func (form *ArgsSource) Get(key string) (string, bool) {
	v := stringsx.FromBytes((*fasthttp.Args)(form).Peek(key))
	return v, v != ""
}

type CtxSource fasthttp.RequestCtx

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *CtxSource) TrySet(value reflect.Value, field *reflect.StructField, key string, opt *kvstruct.Options) (isSet bool, err error) {
	return kvstruct.SetValueByGetter(value, field, form, key, opt)
}

func (form *CtxSource) Get(key string) (string, bool) {
	v := (*fasthttp.RequestCtx)(form).UserValue(key).(string)
	return v, v != ""
}

type HeaderSource fasthttp.RequestHeader

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *HeaderSource) TrySet(value reflect.Value, field *reflect.StructField, key string, opt *kvstruct.Options) (isSet bool, err error) {
	return kvstruct.SetValueByGetter(value, field, form, key, opt)
}

func (form *HeaderSource) Get(key string) (string, bool) {
	v := stringsx.FromBytes((*fasthttp.RequestHeader)(form).Peek(key))
	return v, v != ""
}

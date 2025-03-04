/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/strings"
	"github.com/valyala/fasthttp"
)

type Writer struct {
	fiber.Ctx
}

func (w Writer) Status(code int) {
	w.Ctx.Status(code)
}

func (w Writer) Header() httpi.Header {
	return Header{&w.Ctx.Response().Header}
}

func (w Writer) Write(p []byte) (int, error) {
	return w.Ctx.Write(p)
}

type Header struct {
	*fasthttp.ResponseHeader
}

func (h Header) Add(key, value string) {
	h.ResponseHeader.Add(key, value)
}

func (h Header) Set(key, value string) {
	h.ResponseHeader.Set(key, value)
}

func (h Header) Get(key string) string {
	return strings.BytesToString(h.ResponseHeader.Peek(key))
}

func (h Header) Values(key string) []string {
	byteValues := h.ResponseHeader.PeekAll(key)
	values := make([]string, len(byteValues))
	for i := range byteValues {
		values[i] = strings.BytesToString(byteValues[i])
	}
	return values
}

func (h Header) Range(f func(key, value string)) {
	h.ResponseHeader.VisitAll(func(key, value []byte) {
		f(strings.BytesToString(key), strings.BytesToString(value))
	})
}

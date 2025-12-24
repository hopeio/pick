/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"bufio"
	"context"
	"iter"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/context/reqctx"
	httpx "github.com/hopeio/gox/net/http"
	stringsx "github.com/hopeio/gox/strings"
)

type RequestCtx struct {
	fiber.Ctx
	reqHeader   http.Header
	respHeader  http.Header
	wroteHeader bool
}

func (w RequestCtx) RequestHeader() http.Header {
	if w.reqHeader == nil {
		w.reqHeader = http.Header{}
		w.Ctx.Request().Header.VisitAll(func(k, v []byte) {
			ks := stringsx.FromBytes(k)
			vs := stringsx.FromBytes(v)
			if exists, ok := w.reqHeader[ks]; ok {
				w.reqHeader[ks] = append(exists, vs)
			} else {
				w.reqHeader[ks] = []string{vs}
			}
		})
	}
	return w.reqHeader
}

func (w RequestCtx) RequestContext() context.Context {
	return w.Ctx.Context()
}

func (w RequestCtx) Origin() fiber.Ctx {
	return w.Ctx
}

type Context = reqctx.Context[RequestCtx]

func FromContext(ctx context.Context) (*Context, bool) {
	return reqctx.FromContext[RequestCtx](ctx)
}

func FromRequest(req fiber.Ctx) *Context {
	return reqctx.New[RequestCtx](RequestCtx{Ctx: req})
}

func (w RequestCtx) WriteHeader(code int) {
	w.Ctx.Status(code)
}

func (w RequestCtx) Header() http.Header {
	if w.respHeader == nil {
		w.respHeader = http.Header{}
	}
	return w.respHeader
}

func (w RequestCtx) HeaderX() httpx.Header {
	return ResponseHeader{ResponseHeader: &w.Response().Header}
}

func (w RequestCtx) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		header := &w.Ctx.Response().Header
		for k, v := range w.respHeader {
			for _, vv := range v {
				header.Add(k, vv)
			}
		}
		w.wroteHeader = true
	}
	return w.Ctx.Write(p)
}

func (w RequestCtx) RespondStream(ctx context.Context, dataSource iter.Seq[httpx.WriterToCloser]) {
	w.Ctx.Set(httpx.HeaderTransferEncoding, "chunked")
	w.Ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for data := range dataSource {
			_, err := data.WriteTo(w)
			if err != nil {
				return
			}
			w.Flush()
		}
	})

}

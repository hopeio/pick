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
	iox "github.com/hopeio/gox/io"
	httpx "github.com/hopeio/gox/net/http"
)

type Context struct {
	fiber.Ctx
	reqHeader   http.Header
	respHeader  http.Header
	wroteHeader bool
}

func (w Context) WriteHeader(code int) {
	w.writeHeader()
	w.Ctx.Status(code)
}

func (w Context) writeHeader() {
	if !w.wroteHeader {
		header := &w.Ctx.Response().Header
		for k, v := range w.respHeader {
			for _, vv := range v {
				header.Add(k, vv)
			}
		}
		w.wroteHeader = true
	}
}

func (w Context) Header() http.Header {
	if w.respHeader == nil {
		w.respHeader = http.Header{}
	}
	return w.respHeader
}

func (w Context) HeaderX() httpx.Header {
	return ResponseHeader{ResponseHeader: &w.Response().Header}
}

func (w Context) Write(p []byte) (int, error) {
	w.writeHeader()
	return w.Ctx.Write(p)
}

func (w Context) RespondStream(ctx context.Context, dataSource iter.Seq[iox.WriterToCloser]) {
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

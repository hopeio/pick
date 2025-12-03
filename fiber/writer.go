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

	"github.com/gofiber/fiber/v3"
	httpx "github.com/hopeio/gox/net/http"
)

type Writer struct {
	fiber.Ctx
}

func (w Writer) Status(code int) {
	w.Ctx.Status(code)
}

func (w Writer) Header() httpx.Header {
	return ResponseHeader{ResponseHeader: &w.Response().Header}
}

func (w Writer) Write(p []byte) (int, error) {
	return w.Ctx.Write(p)
}

func (w Writer) RespondStream(ctx context.Context, dataSource iter.Seq[httpx.WriterToCloser]) (int, error) {
	w.Ctx.Set(httpx.HeaderTransferEncoding, "chunked")
	var n, write int64
	var err error
	w.Ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for data := range dataSource {
			write, err = data.WriteTo(w)
			if err != nil {
				return
			}
			n += write
			w.Flush()
		}
	})
	return int(n), nil
}

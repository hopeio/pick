/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	httpi "github.com/hopeio/utils/net/http"
	fiberi "github.com/hopeio/utils/net/http/fiber"
)

type Writer struct {
	fiber.Ctx
}

func (w Writer) Status(code int) {
	w.Ctx.Status(code)
}

func (w Writer) Header() httpi.Header {
	return fiberi.ResponseHeader{ResponseHeader: &w.Response().Header}
}

func (w Writer) Write(p []byte) (int, error) {
	return w.Ctx.Write(p)
}

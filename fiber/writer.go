/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
)

type Writer struct {
	fiber.Ctx
}

func (w Writer) Status(code int) {
	w.Ctx.Status(code)
}

func (w Writer) Set(k, v string) {
	w.Ctx.Set(k, v)
}

func (w Writer) Write(p []byte) (int, error) {
	return w.Ctx.Write(p)
}

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"github.com/gofiber/fiber/v3"
)

type uriSource struct {
	fiber.Ctx
}

func (s uriSource) Get(key string) (string, bool) {
	v := s.Params(key)
	return v, v != ""
}

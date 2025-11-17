/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package middle

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func Log(ctx fiber.Ctx) error {
	log.Println(ctx.Path())
	return ctx.Next()
}

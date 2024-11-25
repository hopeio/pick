/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package middle

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

func Log(ctx fiber.Ctx) error {
	log.Println(ctx.Path())
	return nil
}

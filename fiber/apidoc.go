/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/net/http/apidoc"
	fiberi "github.com/hopeio/gox/net/http/fiber/apidoc"
	apidoc2 "github.com/hopeio/pick/apidoc"
)

func DocList(ctx fiber.Ctx) error {
	modName := ctx.Query("modName")
	if modName == "" {
		modName = "api"
	}
	apidoc2.Openapi(apidoc.Dir, modName)
	return fiberi.DocList(ctx)
}

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/pick"
	"github.com/hopeio/pick/_example/fiber/service"
	fiberi "github.com/hopeio/pick/fiber"
)

func main() {
	app := fiber.New()
	fiberi.Register(app, &service.UserService{}, &service.TestService{})
	pick.OpenApi(":8081")
	log.Println("visit http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

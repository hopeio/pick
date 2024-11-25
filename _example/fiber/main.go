/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/pick/_example/fiber/service"
	_ "github.com/hopeio/pick/_example/fiber/service"
	fiberi "github.com/hopeio/pick/fiber"
	"log"
)

func main() {
	app := fiber.New()
	fiberi.Register(app, &service.UserService{}, &service.TestService{})
	app.Static("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

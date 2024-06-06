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
	fiberi.Start(app, true, &service.UserService{}, &service.TestService{})
	app.Static("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

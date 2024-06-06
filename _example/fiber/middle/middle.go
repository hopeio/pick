package middle

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

func Log(ctx fiber.Ctx) error {
	log.Println(ctx.Path())
	return nil
}

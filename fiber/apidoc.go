package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/net/http/apidoc"
	fiberi "github.com/hopeio/utils/net/http/fasthttp/fiber"
)

func DocList(ctx fiber.Ctx) error {
	modName := ctx.Query("modName")
	if modName == "" {
		modName = "api"
	}
	pick.Markdown(apidoc.ApiDocDir, modName)
	pick.Swagger(apidoc.ApiDocDir, modName)
	return fiberi.DocList(ctx)
}

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
	fiberi "github.com/hopeio/cherry/utils/net/http/fasthttp/fiber"
	"github.com/hopeio/pick"
)

func DocList(ctx fiber.Ctx) {
	modName := ctx.Query("modName")
	if modName == "" {
		modName = "api"
	}
	pick.Markdown(apidoc.ApiDocDir, modName)
	pick.Swagger(apidoc.ApiDocDir, modName)
	fiberi.DocList(ctx)
}

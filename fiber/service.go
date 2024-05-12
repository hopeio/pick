package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/context/fiberctx"
	"github.com/hopeio/pick"
	"reflect"
)

var (
	Svcs             = make([]pick.Service[fiber.Handler], 0)
	FiberContextType = reflect.TypeOf((*fiberctx.Context)(nil))
)

func RegisterService(svc ...pick.Service[fiber.Handler]) {
	Svcs = append(Svcs, svc...)
}

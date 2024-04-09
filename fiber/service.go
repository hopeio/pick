package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/fiber_context"
	"reflect"
)

var (
	Svcs             = make([]pick.Service[fiber.Handler], 0)
	FiberContextType = reflect.TypeOf((*fiber_context.Context)(nil))
)

func RegisterService(svc ...pick.Service[fiber.Handler]) {
	Svcs = append(Svcs, svc...)
}

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/context/fiber_context"
	"github.com/hopeio/pick"
	"reflect"
)

var (
	Svcs             = make([]pick.Service[fiber.Handler], 0)
	FiberContextType = reflect.TypeOf((*fiber_context.Context)(nil))
)

func RegisterService(svc ...pick.Service[fiber.Handler]) {
	Svcs = append(Svcs, svc...)
}

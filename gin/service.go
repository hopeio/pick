package pickgin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/context/gin_context"
	"github.com/hopeio/pick"
	"reflect"
)

var (
	Svcs           = make([]pick.Service[gin.HandlerFunc], 0)
	GinContextType = reflect.TypeOf((*gin_context.Context)(nil))
)

func RegisterService(svc ...pick.Service[gin.HandlerFunc]) {
	Svcs = append(Svcs, svc...)
}

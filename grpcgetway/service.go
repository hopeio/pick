package pickgrpcgetway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pick"
	"reflect"
)

var (
	Svcs        = make([]pick.Service[gin.HandlerFunc], 0)
	ContextType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func RegisterService(svc ...pick.Service[gin.HandlerFunc]) {
	Svcs = append(Svcs, svc...)
}

package pickgin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/context/ginctx"
	"github.com/hopeio/pick"
	"reflect"
)

var (
	Svcs           = make([]pick.Service[gin.HandlerFunc], 0)
	GinContextType = reflect.TypeOf((*ginctx.Context)(nil))
)

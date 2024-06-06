package pickgrpcgetway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/context/ginctx"
	"github.com/hopeio/cherry/utils/log"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/pick"
	"reflect"
)

func Start(engine *gin.Engine, tracing bool, svc ...pick.Service[gin.HandlerFunc]) {
	Svcs = append(Svcs, svc...)
	for _, v := range Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		group := engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, ContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}
			methodType := method.Type
			methodValue := method.Func
			if method.Type.NumIn() < 3 || method.Type.NumOut() != 2 ||
				!methodType.In(1).Implements(ContextType) ||
				!methodType.Out(1).Implements(pick.ErrorType) {
				continue
			}
			methodInfoExport := methodInfo.GetApiInfo()
			in2Type := methodType.In(2)
			group.Handle(methodInfoExport.Method, methodInfoExport.Path, func(ctx *gin.Context) {
				ctxi, s := ginctx.ContextFromRequest(ctx, tracing)
				if s != nil {
					defer s.End()
				}
				in1 := reflect.ValueOf(ctxi).Interface().(context.Context)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, reflect.ValueOf(in1), in2})
				httpi.ResWriteReflect(ctx.Writer, ctxi.TraceID, result)
			})
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}

	pick.Registered(Svcs)
}

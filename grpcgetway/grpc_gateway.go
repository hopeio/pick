package pickgrpcgetway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/log"
	"reflect"
)

var (
	ContextType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func Register(engine *gin.Engine, tracing bool, svcs ...pick.Service[gin.HandlerFunc]) {
	for _, v := range svcs {
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
				ctxi := ginctx.FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in1 := reflect.ValueOf(ctxi).Interface().(context.Context)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, reflect.ValueOf(in1), in2})
				pick.ResWriteReflect(ctx.Writer, ctxi.TraceID(), result)
			})
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}

	pick.Registered()
}

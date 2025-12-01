package pickgin

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/gin/binding"
	"github.com/hopeio/pick"
)

func RegisterGrpcService(engine *gin.Engine, svcs ...pick.Service[gin.HandlerFunc]) {
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*pick.ApiDocInfo
		group := engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, GinContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}
			methodType := method.Type
			if !methodType.In(1).Implements(pick.ContextType) {
				continue
			}
			methodValue := method.Func
			in2Type := methodType.In(2).Elem()
			methodInfoExport := methodInfo.Export()
			handler := func(ctx *gin.Context) {
				ctxi := httpctx.FromRequest(httpctx.RequestCtx{Request: ctx.Request, ResponseWriter: ctx.Writer})
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				err := binding.Bind(ctx, in2.Interface())
				if err != nil {
					ctx.JSON(http.StatusBadRequest, errors.InvalidArgument.Msg(err.Error()))
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				params[1] = reflect.ValueOf(ctxi.Wrapper())
				params[2] = in2
				result := methodValue.Call(params)
				pick.Respond(ctx, Writer{ctx}, ctxi.TraceID(), result)
			}
			for _, url := range methodInfoExport.Routes {
				group.Handle(url.Method, url.Path[len(preUrl):], handler)
			}
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

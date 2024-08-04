package pickgin

import (
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/errors/errcode"
	gin2 "github.com/hopeio/utils/net/http/gin"

	"github.com/hopeio/utils/net/http/apidoc"
	binding "github.com/hopeio/utils/net/http/gin/binding"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	GinContextType = reflect.TypeOf((*ginctx.Context)(nil))
)

func Register(engine *gin.Engine, svcs ...pick.Service[gin.HandlerFunc]) {
	gin2.OpenApi(engine, "", "")
	openApiLog()
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
			methodValue := method.Func
			in2Type := methodType.In(2)
			methodInfoExport := methodInfo.GetApiInfo()
			group.Handle(methodInfoExport.Method, methodInfoExport.Path[len(preUrl):], func(ctx *gin.Context) {
				ctxi := ginctx.FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				err := binding.Bind(ctx, in2.Interface())
				if err != nil {
					ctx.JSON(http.StatusBadRequest, errcode.InvalidArgument.Msg(err.Error()))
					return
				}
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				pick.ResWriteReflect(ctx.Writer, ctxi.TraceID(), result)
			})
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

func openApiLog() {
	pick.Log(http.MethodGet, apidoc.UriPrefix, "api文档列表")
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/swagger/*file", "swagger文档")
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/markdown/*file", "markdown文档")
}

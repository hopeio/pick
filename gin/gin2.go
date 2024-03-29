package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"github.com/hopeio/tiga/protobuf/errorcode"
	"github.com/hopeio/tiga/utils/net/http/api/apidoc"
	gin_build "github.com/hopeio/tiga/utils/net/http/gin"
	"github.com/hopeio/tiga/utils/net/http/gin/handler"
	"github.com/hopeio/tiga/utils/net/http/request"
	"log"
	"net/http"
	"reflect"
)

func genApi(middlewareHandler func(preUrl string, middleware []http.HandlerFunc), handle func(method, path string, in2Type reflect.Type, methodValue, value reflect.Value)) {
	for _, v := range pick.Svcs {
		describe, preUrl, middleware := v.Service()
		middlewareHandler(preUrl, middleware)
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, pick.ClaimsType)
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
			handle(methodInfoExport.Method, methodInfoExport.Path, in2Type, methodValue, value)
			infos = append(infos, &pick.ApiDocInfo{methodInfo, method.Type})
		}
		pick.GroupApiInfos = append(pick.GroupApiInfos, &pick.GroupApiInfo{describe, infos})
	}

	pick.Registered()
}

func RegisterGinAPI(engine *gin.Engine, genDoc bool, modName string, tracing bool) {
	genApi(func(preUrl string, middleware []http.HandlerFunc) {
		engine.Group(preUrl, handler.Converts(middleware)...)
	},
		func(method, path string, in2Type reflect.Type, methodValue, value reflect.Value) {
			engine.Handle(method, path, func(ctx *gin.Context) {
				ctxi, span := http_context.ContextFromRequestResponse(ctx.Request, ctx.Writer, tracing)
				if span != nil {
					defer span.End()
				}
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				err := gin_build.Bind(ctx, in2.Interface())
				if err != nil {
					ctx.JSON(http.StatusBadRequest, errorcode.InvalidArgument.Message(request.Error(err)))
					return
				}
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				pick.ResHandler(ctxi, ctx.Writer, result)
			})
		})
	if genDoc {
		pick.GenApiDoc(modName)
		gin_build.OpenApi(engine, apidoc.FilePath)
	}

}

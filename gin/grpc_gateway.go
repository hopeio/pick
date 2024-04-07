package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"github.com/hopeio/tiga/utils/log"
	"github.com/hopeio/tiga/utils/net/http/api/apidoc"
	gin_build "github.com/hopeio/tiga/utils/net/http/gin"
	"github.com/hopeio/tiga/utils/net/http/gin/handler"
	"reflect"
)

func RegisterGrpcGateway(engine *gin.Engine, genDoc bool, modName string, tracing bool) {

	for _, v := range pick.Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		group := engine.Group(preUrl, handler.Converts(middleware)...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, pick.HttpContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}
			methodType := method.Type
			methodValue := method.Func
			if method.Type.NumIn() < 3 || method.Type.NumOut() != 2 ||
				!methodType.In(1).Implements(pick.ContextType) ||
				!methodType.Out(1).Implements(pick.ErrorType) {
				continue
			}

			in2Type := methodType.In(2)
			group.Handle(methodInfo.method, methodInfo.path, func(ctx *gin.Context) {
				ctxi, s := http_context.ContextFromRequestResponse(ctx.Request, ctx.Writer, tracing)
				if s != nil {
					defer s.End()
				}
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				pick.ResHandler(ctxi, ctx.Writer, result)
			})
			infos = append(infos, &pick.ApiDocInfo{methodInfo, method.Type})
		}
		pick.GroupApiInfos = append(pick.GroupApiInfos, &pick.GroupApiInfo{describe, infos})
	}
	if genDoc {
		pick.GenApiDoc(modName)
		gin_build.OpenApi(engine, apidoc.FilePath)
	}

}

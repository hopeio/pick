package pickgrpcgetway

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pick"
	"github.com/hopeio/pick/gin"
	"github.com/hopeio/tiga/context/gin_context"
	"github.com/hopeio/tiga/utils/log"
	"github.com/hopeio/tiga/utils/net/http/api/apidoc"
	gin_build "github.com/hopeio/tiga/utils/net/http/gin"
	"reflect"
)

func Start(engine *gin.Engine, genDoc bool, modName string, tracing bool) {

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
			methodInfo := pick.GetMethodInfo(&method, preUrl, pickgin.GinContextType)
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
				ctxi, s := gin_context.ContextFromRequest(ctx, tracing)
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
	pick.Registered(Svcs)
}

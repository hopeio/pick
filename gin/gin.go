package pickgin

import (
	"github.com/hopeio/cherry/context/gin_context"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/net/http/request"
	"github.com/hopeio/pick"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
	gin_build "github.com/hopeio/cherry/utils/net/http/gin"
)

// 虽然我写的路由比httprouter更强大(没有map,lru cache)，但是还是选择用gin,理由是gin也用同样的方式改造了路由

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
				ctxi, span := gin_context.ContextFromRequest(ctx, tracing)
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
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.GroupApiInfos = append(pick.GroupApiInfos, &pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	if genDoc {
		pick.GenApiDoc(modName)
		gin_build.OpenApi(engine, apidoc.FilePath)
	}
	pick.Registered(Svcs)
}

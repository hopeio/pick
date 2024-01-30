package gin

import (
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/tiga/utils/log"
	"github.com/hopeio/tiga/utils/net/http/api/apidoc"
	"github.com/hopeio/tiga/utils/net/http/gin/handler"
)

// Deprecated:这种方法不推荐使用了，目前就两种定义api的方式，一种grpc-gateway，一种pick自定义
// 该方法适用于不使用grpc-gateway的情况，只用该方法定义api
func GrpcServiceToRestfulApi(engine *gin.Engine, genApi bool, modName string, tracing bool) {
	httpMethods := []string{http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodConnect, http.MethodHead, http.MethodTrace}
	doc := apidoc.GetDoc(filepath.Join(apidoc.FilePath+modName, modName+apidoc.EXT))
	methods := make(map[string]struct{})
	for _, v := range pick.Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		group := engine.Group(preUrl, handler.Converts(middleware)...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodType := method.Type
			methodValue := method.Func
			if method.Type.NumIn() < 3 || method.Type.NumOut() != 2 ||
				!methodType.In(1).Implements(pick.ContextType) ||
				!methodType.Out(1).Implements(pick.ErrorType) {
				continue
			}

			methodInfo := new(pick.ApiInfo)
			methodInfo.Title = describe
			methodInfo.Middleware = middleware
			methodInfo.Method, methodInfo.Path, methodInfo.Version = pick.ParseMethodName(method.Name, httpMethods)
			methodInfo.Path = "/api/v" + strconv.Itoa(methodInfo.Version) + "/" + methodInfo.Path

			in2Type := methodType.In(2)
			group.Handle(methodInfo.Method, methodInfo.Path, func(ctx *gin.Context) {
				ctxi, s := http_context.ContextFromRequest(ctx.Request, tracing)
				if s != nil {
					defer s.End()
				}
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				pick.ResHandler(ctxi, ctx.Writer, result)
			})
			methods[methodInfo.Method] = struct{}{}
			if genApi {
				methodInfo.GetApiInfo().Swagger(doc, value.Method(j).Type(), describe, value.Type().Name())
			}
		}

	}
	if genApi {
		apidoc.WriteToFile(apidoc.FilePath, modName)
	}

}

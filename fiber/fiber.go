package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/fiberctx"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/errors/errcode"
	"github.com/hopeio/utils/net/http/apidoc"
	fiberi "github.com/hopeio/utils/net/http/fasthttp/fiber/apidoc"
	"github.com/hopeio/utils/net/http/fasthttp/fiber/binding"
	"net/http"
	"reflect"

	"github.com/hopeio/utils/log"
)

var (
	FiberContextType = reflect.TypeOf((*fiberctx.Context)(nil))
)

// 复用pick service，不支持单个接口的中间件
func Register(engine *fiber.App, svcs ...pick.Service[fiber.Handler]) {
	openApi(engine)
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
			methodInfo := pick.GetMethodInfo(&method, preUrl, FiberContextType)
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
			group.Add([]string{methodInfoExport.Method}, methodInfoExport.Path[len(preUrl):], func(ctx fiber.Ctx) error {
				ctxi := fiberctx.FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				if err := binding.Bind(ctx, in2.Interface()); err != nil {
					return ctx.Status(http.StatusBadRequest).JSON(errcode.InvalidArgument.ErrRep())
				}
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				return ResWriterReflect(ctx, ctxi.TraceID(), result)
			})
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}

	pick.Registered()
}

func openApi(mux *fiber.App) {
	pick.Log(http.MethodGet, apidoc.UriPrefix, "api文档列表")
	mux.Get(apidoc.UriPrefix, DocList)
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/swagger/*file", "swagger文档")
	mux.Get(apidoc.UriPrefix+"/swagger/*file", fiberi.Swagger)
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/markdown/*file", "markdown文档")
	mux.Get(apidoc.UriPrefix+"/markdown", fiberi.Markdown)
}

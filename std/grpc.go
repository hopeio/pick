package std

import (
	"encoding/json"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/pick"
	apidoc2 "github.com/hopeio/pick/apidoc"
	"github.com/hopeio/gox/errors/errcode"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/binding"
	"github.com/hopeio/gox/net/http/consts"
	"github.com/hopeio/gox/unsafe"
	"net/http"
	"reflect"
)

func RegisterGrpcService(engine *http.ServeMux, svcs ...pick.Service[Middleware]) {
	openApi(engine)
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*apidoc2.ApiDocInfo

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo[Middleware](&method, preUrl, HttpContextType)
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
			handler := func(w http.ResponseWriter, r *http.Request) {
				ctxi := httpctx.FromRequest(httpctx.RequestCtx{Request: r, Response: w})
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				err := binding.Bind(r, in2.Interface())
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set(consts.HeaderContentType, consts.ContentTypeJsonUtf8)
					json.NewEncoder(w).Encode(errcode.InvalidArgument.Msg(err.Error()))
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				params[1] = reflect.ValueOf(ctxi.Wrapper())
				params[2] = in2
				result := methodValue.Call(params)
				pick.Response(Writer{w}, ctxi.TraceID(), result)
			}
			for _, url := range methodInfoExport.Routes {
				engine.Handle(url.Method+" "+url.Path[len(preUrl):], Chain(handler, append(middleware, unsafe.CastSlice[Middleware](methodInfoExport.Middlewares)...)...))
			}
			methodInfo.Log()
			infos = append(infos, &apidoc2.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		apidoc2.RegisterApiInfo(&apidoc2.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

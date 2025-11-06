package std

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	http2 "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	"github.com/hopeio/gox/net/http/binding"
	"github.com/hopeio/gox/unsafe"
	"github.com/hopeio/pick"
	apidoc2 "github.com/hopeio/pick/apidoc"
)

var (
	HttpContextType = reflect.TypeOf((*httpctx.Context)(nil))
)

func Register(engine *http.ServeMux, svcs ...pick.Service[Middleware]) {
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
			methodValue := method.Func
			in2Type := methodType.In(2).Elem()
			methodInfoExport := methodInfo.Export()
			httpContext := methodType.In(1).ConvertibleTo(HttpContextType)
			handler := func(w http.ResponseWriter, r *http.Request) {
				ctxi := httpctx.FromRequest(httpctx.RequestCtx{Request: r, ResponseWriter: w})
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				err := binding.Bind(r, in2.Interface())
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set(http2.HeaderContentType, http2.ContentTypeJsonUtf8)
					json.NewEncoder(w).Encode(errors.InvalidArgument.Msg(err.Error()))
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				if httpContext {
					params[1] = reflect.ValueOf(ctxi)
				} else {
					params[1] = reflect.ValueOf(ctxi.Wrapper())
				}
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

func openApi(mux *http.ServeMux) {
	mux.HandleFunc(apidoc.UriPrefix, apidoc2.DocList)
	pick.Log(http.MethodGet, apidoc.UriPrefix, "apidoc list")
	mux.HandleFunc(apidoc.UriPrefix+"/openapi/*file", apidoc.OpenApi)
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/openapi/*file", "openapi")
}

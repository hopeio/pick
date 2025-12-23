package pickstd

import (
	"net/http"
	"reflect"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/pick"
)

func RegisterGrpcService(engine *http.ServeMux, svcs ...pick.Service[Middleware]) {
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*pick.ApiDocInfo

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, HttpContextType)
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
				ctxi := httpctx.FromRequest(httpctx.RequestCtx{Request: r, ResponseWriter: w})
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				reqv := in2.Interface()
				err := httpx.Bind(r, reqv)
				if err != nil {
					pick.RespondError(ctxi.Base(), httpx.ResponseWriterWrapper{w}, errors.InvalidArgument.Msg(err.Error()), ctxi.TraceID())
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				params[1] = reflect.ValueOf(ctxi.Wrapper())
				params[2] = in2
				result := methodValue.Call(params)
				pick.Respond(ctxi.Base(), httpx.ResponseWriterWrapper{w}, ctxi.TraceID(), result)
			}
			for _, url := range methodInfoExport.Routes {
				engine.Handle(url.Method+" "+url.Path[len(preUrl):], httpx.UseMiddleware(http.HandlerFunc(handler), middleware...))
			}
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

package pickstd

import (
	"context"
	"net/http"
	"reflect"

	errorsx "github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/pick"
)

var (
	HttpContextType = reflect.TypeOf((*Context)(nil))
)

func Register(engine *http.ServeMux, svcs ...pick.Service[Middleware]) {
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*pick.ApiDocInfo

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, HttpContextType, reflect.ValueOf(&Context{}))
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
				in2 := reflect.New(in2Type)
				err := httpx.Bind(r, in2.Interface())
				if err != nil {
					pick.RespondError(r.Context(), w, errorsx.InvalidArgument.Msg(err.Error()))
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				if httpContext {
					params[1] = reflect.ValueOf(&Context{r, w})
				} else {
					r.WithContext(context.WithValue(r.Context(), httpx.RequestCtxKey, &Context{r, w}))
					params[1] = reflect.ValueOf(r.Context())
				}
				params[2] = in2
				result := methodValue.Call(params)
				pick.Respond(r.Context(), w, result)
			}
			for _, url := range methodInfoExport.Routes {
				engine.Handle(url.Method+" "+url.Path, httpx.UseMiddleware(http.HandlerFunc(handler), middleware...))
			}
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

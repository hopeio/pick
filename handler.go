package pick

import (
	"encoding/json"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/reflect/mtos"
	"net/http"
	"reflect"
)

var (
	HttpContextType = reflect.TypeOf((*httpctx.Context)(nil))
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

func CommonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		ctxi := httpctx.FromRequest(httpctx.RequestCtx{
			Request:  req,
			Response: w,
		})
		defer ctxi.RootSpan().End()
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).ConvertibleTo(HttpContextType) {
				params[i] = reflect.ValueOf(ctxi)
			} else {
				params[i] = reflect.New(handleTyp.In(i).Elem())
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					if ps != nil {
						pathParam := *ps
						if len(pathParam) > 0 {
							for i := range pathParam {
								src.Set(pathParam[i].Key, pathParam[i].Value)
							}
						}
					}
					mtos.Decode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		ResWriteReflect(w, ctxi.TraceID(), result)
	}
}

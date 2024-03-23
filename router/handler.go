package router

import (
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"github.com/hopeio/tiga/utils/encoding/json"
	"github.com/hopeio/tiga/utils/net/http/request/binding"
	"net/http"
	"reflect"
)

func commonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params, tracing bool) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		ctxi, s := http_context.ContextFromRequestResponse(req, w, tracing)
		if s != nil {
			defer s.End()
		}
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).Implements(pick.ClaimsType) {
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
					binding.PickDecode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		pick.ResHandler(ctxi, w, result)
	}
}

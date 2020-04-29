package pick

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/liov/pick/response"
	"github.com/liov/pick/schema"
)

type Handlers []http.Handler

func (hs Handlers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range hs {
		handler.ServeHTTP(w, req)
	}
}

type HandlerFuncs []http.HandlerFunc

func (hs HandlerFuncs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range hs {
		handler(w, req)
	}
}

func (hs HandlerFuncs) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, handler := range hs {
			handler(w, req)
		}
	}
}

func (hs *HandlerFuncs) Add(handler http.HandlerFunc) {
	*hs = append(*hs, handler)
}

func commonHandler(w http.ResponseWriter, req *http.Request, handle reflect.Value) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		for i := 0; i < handleNumIn; i++ {
			params[i] = reflect.New(handleTyp.In(i).Elem())
			if handleTyp.In(i).Implements(contextType) {
				sess := params[i].Interface().(Session)
				sess.Parse(req)
			} else {
				if req.URL.RawQuery != "" {
					src := req.URL.Query()
					decoder.PickDecode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		if !result[1].IsNil() {
			json.NewEncoder(w).Encode(result[1].Interface())
			return
		}
		if info, ok := result[0].Interface().(*response.File); ok {
			header := w.Header()
			header.Set("Content-Type", "application/octet-stream")
			header.Set("Content-Disposition", "attachment;filename="+info.Name)
			io.Copy(w, info.File)
			if flusher, canFlush := w.(http.Flusher); canFlush {
				flusher.Flush()
			}
			info.File.Close()
			return
		}
		json.NewEncoder(w).Encode(response.ResData{
			Code:    0,
			Message: "success",
			Details: result[0].Interface(),
		})
	}
}

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
}

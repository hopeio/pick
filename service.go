package pick

import (
	"context"
	"encoding/json"
	"github.com/hopeio/tiga/context/http_context"
	"io"
	"net/http"
	"reflect"

	"github.com/hopeio/tiga/protobuf/errorcode"
	httpi "github.com/hopeio/tiga/utils/net/http"
	http_fs "github.com/hopeio/tiga/utils/net/http/fs"
)

type ParseToHttpResponse interface {
	Parse() ([]byte, error)
}

var (
	Svcs            = make([]Service, 0)
	isRegistered    = false
	HttpContextType = reflect.TypeOf((*http_context.Context)(nil))
	ContextType     = reflect.TypeOf((*context.Context)(nil)).Elem()
	ErrorType       = reflect.TypeOf((*error)(nil)).Elem()
)

type Service interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []http.HandlerFunc)
}

func RegisterService(svc ...Service) {
	Svcs = append(Svcs, svc...)
}

func Registered() {
	isRegistered = true
	Svcs = nil
	GroupApiInfos = nil
}

func Api(f func()) {
	if !isRegistered {
		f()
	}
}

// 兼容有返回值和无返回值的写法
func Api2(f func() any) {
	if !isRegistered {
		panic(f())
	}
}

func ResHandler(c *http_context.Context, w http.ResponseWriter, result []reflect.Value) {
	if !result[1].IsNil() {
		err := errorcode.ErrHandle(result[1].Interface())
		c.HandleError(err)
		json.NewEncoder(w).Encode(err)
		return
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := w.Header()
		header.Set(httpi.HeaderContentType, httpi.ContentBinaryHeaderValue)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(w, info.File)
		if flusher, canFlush := w.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	json.NewEncoder(w).Encode(httpi.ResAnyData{
		Message: "OK",
		Details: result[0].Interface(),
	})
}

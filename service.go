package pick

import (
	"encoding/json"
	contexti "github.com/hopeio/cherry/context"
	"io"
	"net/http"
	"reflect"

	"github.com/hopeio/cherry/protobuf/errorcode"
	httpi "github.com/hopeio/cherry/utils/net/http"
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
)

type ParseToHttpResponse interface {
	Parse() ([]byte, error)
}

var (
	isRegistered = false
	ErrorType    = reflect.TypeOf((*error)(nil)).Elem()
)

type Service[T any] interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []T)
}

func Registered[T any](svcs []Service[T]) {
	isRegistered = true
	svcs = nil
	GroupApiInfos = nil
}

func Api(f func()) {
	if !isRegistered {
		f()
	}
}

func ResHandler[T any](c *contexti.RequestContext[T], w http.ResponseWriter, result []reflect.Value) {
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

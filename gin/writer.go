package pickgin

import (
	"github.com/gin-gonic/gin"
	httpi "github.com/hopeio/utils/net/http"
)

type Writer struct {
	*gin.Context
}

func (w Writer) Status(code int) {
	w.Context.Status(code)
}

func (w Writer) Set(k, v string) {
	w.Context.Header(k, v)
}

func (w Writer) Write(p []byte) (int, error) {
	return w.Context.Writer.Write(p)
}

type Response interface {
	func(*gin.Context)
}

func ResponseHook(ctx *gin.Context) func(any) (bool, error) {
	return func(data any) (bool, error) {

		if info, ok := data.(httpi.IHttpResponseTo); ok {
			_, err := info.Response(ctx.Writer)
			if err != nil {
				return false, err
			}
			return true, err
		}
		return false, nil
	}
}

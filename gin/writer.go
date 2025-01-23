package pickgin

import "github.com/gin-gonic/gin"

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

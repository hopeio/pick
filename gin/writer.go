package pickgin

import (
	"context"
	"iter"

	"github.com/gin-gonic/gin"
	httpx "github.com/hopeio/gox/net/http"
)

type Writer struct {
	*gin.Context
}

func (w Writer) Status(code int) {
	w.Context.Status(code)
}

func (w Writer) Header() httpx.Header {
	return httpx.HttpHeader(w.Writer.Header())
}
func (w Writer) Write(p []byte) (int, error) {
	return w.Context.Writer.Write(p)
}

func (w Writer) RespondStream(ctx context.Context, seq iter.Seq[httpx.WriterToCloser]) (int, error) {
	return httpx.RespondStream(ctx, w.Writer, seq)
}

func (w Writer) RecordResponse(contentType string, raw []byte, res any) {
	if w, ok := w.Writer.(httpx.Unwrapper); ok {
		if recorder, ok := w.Unwrap().(httpx.ResponseRecorder); ok {
			recorder.RecordResponse(contentType, raw, res)
		}
	}
}

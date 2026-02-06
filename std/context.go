package pickstd

import "net/http"

type Context struct {
	*http.Request
	http.ResponseWriter
}

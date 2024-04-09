package pickrouter

import (
	"github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"net/http"
	"reflect"
)

var (
	Svcs            = make([]pick.Service[http.HandlerFunc], 0)
	HttpContextType = reflect.TypeOf((*http_context.Context)(nil))
)

func RegisterService(svc ...pick.Service[http.HandlerFunc]) {
	Svcs = append(Svcs, svc...)
}

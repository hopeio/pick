package pick

import (
	"net/http"
	"reflect"
)

type Session interface {
	Parse(*http.Request) error
}

var contextType = reflect.TypeOf((*Session)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

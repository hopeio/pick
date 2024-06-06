package pick

import (
	"reflect"
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

func Registered() {
	isRegistered = true
}

func Api(f func()) {
	if !isRegistered {
		f()
	}
}

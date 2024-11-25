/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

type ParseToHttpResponse interface {
	Parse() ([]byte, error)
}

var (
	isRegistered = false
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

var prefix string

func HandlerPrefix(p string) {
	prefix = p
}

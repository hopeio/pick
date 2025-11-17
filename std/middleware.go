package std

import (
	httpx "github.com/hopeio/gox/net/http"
)

type Middleware = httpx.Middleware

// UseMiddleware applies middlewares to a http.HandlerFunc
var UseMiddleware = httpx.UseMiddleware

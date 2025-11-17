/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"log"

	pickstd "github.com/hopeio/pick/std"
)

func Log(ctx *pickstd.MiddlewareContext) {
	log.Println("Log", ctx.Request.Method, ctx.Request.RequestURI)
	ctx.Next()
	log.Println("Log End", ctx.Request.Method, ctx.Request.RequestURI)
}

func Log2(ctx *pickstd.MiddlewareContext) {
	log.Println("Log2", ctx.Request.Method, ctx.Request.RequestURI)
	ctx.Next()
	log.Println("Log2 End", ctx.Request.Method, ctx.Request.RequestURI)
}

func Log3(ctx *pickstd.MiddlewareContext) {
	log.Println("Log3", ctx.Request.Method, ctx.Request.RequestURI)
	ctx.Next()
	log.Println("Log3 End", ctx.Request.Method, ctx.Request.RequestURI)
}

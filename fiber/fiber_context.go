/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/reqctx"
	httpx "github.com/hopeio/gox/net/http"
)

type RequestCtx struct {
	fiber.Ctx
}

func (ctx RequestCtx) RequestHeader() httpx.Header {
	return RequestHeader{RequestHeader: &ctx.Request().Header}
}

func (ctx RequestCtx) RequestContext() context.Context {
	return ctx.Ctx.Context()
}

func (ctx RequestCtx) Origin() fiber.Ctx {
	return ctx.Ctx
}

type Context = reqctx.Context[RequestCtx]

func FromContext(ctx context.Context) (*Context, bool) {
	return reqctx.FromContext[RequestCtx](ctx)
}

func FromRequest(req fiber.Ctx) *Context {
	return reqctx.New[RequestCtx](RequestCtx{req})
}

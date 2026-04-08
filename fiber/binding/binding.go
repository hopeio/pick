/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"context"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v3"
	iox "github.com/hopeio/gox/io"
	"github.com/hopeio/gox/mapstruct"
	httpx "github.com/hopeio/gox/net/http"
	stringsx "github.com/hopeio/gox/strings"
)

func Bind(c fiber.Ctx, obj interface{}) error {
	return httpx.CommonBind(RequestSource{c}, obj)
}

type RequestSource struct {
	fiber.Ctx
}

func (s RequestSource) Uri() mapstruct.Getter {
	return uriSource{s.Ctx}
}

func (s RequestSource) Query() mapstruct.ValuesGetter {
	return (*ArgsSource)(s.Request().URI().QueryArgs())
}

func (s RequestSource) Header() mapstruct.ValuesGetter {
	return (*HeaderSource)(&s.Request().Header)
}

func (s RequestSource) Body() (context.Context, string, io.ReadCloser) {
	if s.Method() == http.MethodGet {
		return s.Context(), "", nil
	}
	contentType := stringsx.FromBytes(s.Request().Header.ContentType())
	req := s.Ctx.Request()
	if req.IsBodyStream() {
		return s.Context(), contentType, iox.WrapReader(req.BodyStream(), req.CloseBodyStream)
	}
	return s.Context(), contentType, iox.RawBytes(req.Body())
}

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v3"
	iox "github.com/hopeio/gox/io"
	"github.com/hopeio/gox/kvstruct"
	httpx "github.com/hopeio/gox/net/http"
	stringsx "github.com/hopeio/gox/strings"
)

func Bind(c fiber.Ctx, obj interface{}) error {
	return httpx.CommonBind(RequestSource{c}, obj)
}

type RequestSource struct {
	fiber.Ctx
}

func (s RequestSource) Uri() kvstruct.Setter {
	return uriSource{s.Ctx}
}

func (s RequestSource) Query() kvstruct.Setter {
	return (*ArgsSource)(s.Request().URI().QueryArgs())
}

func (s RequestSource) Header() kvstruct.Setter {
	return (*HeaderSource)(&s.Request().Header)
}

func (s RequestSource) Form() kvstruct.Setter {
	contentType := stringsx.FromBytes(s.Request().Header.ContentType())
	if strings.HasPrefix(contentType, httpx.ContentTypeForm) {
		vs, err := url.ParseQuery(stringsx.FromBytes(s.Ctx.Body()))
		if err != nil {
			return nil
		}
		return kvstruct.KVsSource(vs)
	}
	if strings.HasPrefix(contentType, httpx.ContentTypeMultipart) {
		multipartForm, err := s.Ctx.MultipartForm()
		if err != nil {
			return nil
		}
		return (*httpx.MultipartSource)(multipartForm)
	}
	return nil
}

func (s RequestSource) Body() (string, io.ReadCloser) {
	if s.Method() == http.MethodGet {
		return "", nil
	}
	if s.Is(httpx.ContentTypeMultipart) || s.Is(httpx.ContentTypeForm) {
		return "", nil
	}
	contentType := stringsx.FromBytes(s.Request().Header.ContentType())
	req := s.Ctx.Request()
	if req.IsBodyStream() {
		return contentType, iox.WrapReader(req.BodyStream(), req.CloseBodyStream)
	}
	return contentType, iox.RawBytes(req.Body())
}

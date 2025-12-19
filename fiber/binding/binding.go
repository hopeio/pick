/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/kvstruct"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/binding"
	stringsx "github.com/hopeio/gox/strings"
)

func Bind(c fiber.Ctx, obj interface{}) error {
	return binding.CommonBind(RequestSource{c}, obj)
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

func (s RequestSource) MultipartForm() kvstruct.Setter {
	contentType := stringsx.FromBytes(s.Request().Header.Peek(httpx.HeaderContentType))
	if strings.HasPrefix(contentType, httpx.ContentTypeMultipart) {
		multipartForm, err := s.Ctx.MultipartForm()
		if err != nil {
			return nil
		}
		return (*binding.MultipartSource)(multipartForm)
	}
	return nil
}

func (s RequestSource) BodyBind(obj any) error {
	if s.Method() == http.MethodGet {
		return nil
	}
	if s.Is(httpx.ContentTypeMultipart) {
		return nil
	}
	return binding.BodyUnmarshaller(stringsx.FromBytes(s.Request().Header.Peek(httpx.HeaderContentType)), s.Body(), obj)
}

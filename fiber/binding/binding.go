/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/mtos"
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

func (s RequestSource) Uri() mtos.Setter {
	return uriSource{s.Ctx}
}

func (s RequestSource) Query() mtos.Setter {
	return (*ArgsSource)(s.Request().URI().QueryArgs())
}

func (s RequestSource) Header() mtos.Setter {
	return (*HeaderSource)(&s.Request().Header)
}

func (s RequestSource) Form() mtos.Setter {
	contentType := stringsx.FromBytes(s.Request().Header.Peek(httpx.HeaderContentType))
	if contentType == httpx.ContentTypeForm {
		return (*ArgsSource)(s.Request().PostArgs())
	}
	if contentType == httpx.ContentTypeMultipart {
		multipartForm, err := s.MultipartForm()
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
	return binding.BodyUnmarshaller(s.Body(), obj)
}

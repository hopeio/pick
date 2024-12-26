/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/fiberctx"
	"github.com/hopeio/pick"
	apidoc2 "github.com/hopeio/pick/apidoc"
	"github.com/hopeio/utils/errors/errcode"
	"github.com/hopeio/utils/net/http/apidoc"
	fiberi "github.com/hopeio/utils/net/http/fiber/apidoc"
	"github.com/hopeio/utils/net/http/fiber/binding"
	"net/http"
	"reflect"

	"github.com/hopeio/utils/log"
)

var (
	FiberContextType = reflect.TypeOf((*fiberctx.Context)(nil))
)

// 复用pick service，不支持单个接口的中间件
func Register(engine *fiber.App, svcs ...pick.Service[fiber.Handler]) {
	openApi(engine)
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*apidoc2.ApiDocInfo
		group := engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, FiberContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}

			methodType := method.Type
			methodValue := method.Func
			in2Type := methodType.In(2).Elem()
			methodInfoExport := methodInfo.Export()
			group.Add([]string{methodInfoExport.Method}, methodInfoExport.Path[len(preUrl):], func(ctx fiber.Ctx) error {
				ctxi := fiberctx.FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				if err := binding.Bind(ctx, in2.Interface()); err != nil {
					return ctx.Status(http.StatusBadRequest).JSON(errcode.InvalidArgument.ErrRep())
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				if methodType.In(1).ConvertibleTo(FiberContextType) {
					params[1] = reflect.ValueOf(ctxi)
				} else {
					params[1] = reflect.ValueOf(ctxi.Wrapper())
				}
				params[2] = in2
				result := methodValue.Call(params)
				return ResWriterReflect(ctx, ctxi.TraceID(), result)
			})
			methodInfo.Log()
			infos = append(infos, &apidoc2.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		apidoc2.RegisterApiInfo(&apidoc2.GroupApiInfo{Describe: describe, Infos: infos})
	}

	pick.Registered()
}

func openApi(mux *fiber.App) {
	pick.Log(http.MethodGet, apidoc.UriPrefix, "api文档列表")
	mux.Get(apidoc.UriPrefix, DocList)
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/openapi/*file", "openapi文档")
	mux.Get(apidoc.UriPrefix+"/openapi/*file", fiberi.Swagger)
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/markdown/*file", "markdown文档")
	mux.Get(apidoc.UriPrefix+"/markdown", fiberi.Markdown)
}

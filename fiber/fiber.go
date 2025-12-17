/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/pick"
	"github.com/hopeio/pick/fiber/binding"

	"github.com/hopeio/gox/log"
)

var (
	FiberContextType = reflect.TypeOf((*Context)(nil))
)

// 复用pick service，不支持单个接口的中间件
func Register(engine *fiber.App, svcs ...pick.Service[fiber.Handler]) {
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("service must be a pointer")
		}
		var infos []*pick.ApiDocInfo
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
			fiberContext := methodType.In(1).ConvertibleTo(FiberContextType)
			handler := func(ctx fiber.Ctx) error {
				ctxi := FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				if err := binding.Bind(ctx, in2.Interface()); err != nil {
					return pick.RespondError(ctx.Context(), Writer{ctx}, errors.InvalidArgument.Msg(err.Error()), ctxi.TraceID())
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				if fiberContext {
					params[1] = reflect.ValueOf(ctxi)
				} else {
					params[1] = reflect.ValueOf(ctxi.Wrapper())
				}
				params[2] = in2
				result := methodValue.Call(params)
				return pick.Respond(ctx.Context(), Writer{ctx}, ctxi.TraceID(), result)
			}
			for _, url := range methodInfoExport.Routes {
				group.Add([]string{url.Method}, url.Path[len(preUrl):], handler)
			}
			methodInfo.Log()
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		pick.RegisterApiInfo(&pick.GroupApiInfo{Describe: describe, Infos: infos})
	}

	pick.Registered()
}

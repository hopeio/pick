/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickgin

import (
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/pick"
	apidoc2 "github.com/hopeio/pick/apidoc"
	"github.com/hopeio/utils/errors/errcode"
	gin2 "github.com/hopeio/utils/net/http/gin"

	"github.com/hopeio/utils/net/http/apidoc"
	binding "github.com/hopeio/utils/net/http/gin/binding"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	GinContextType = reflect.TypeOf((*ginctx.Context)(nil))
)

func Register(engine *gin.Engine, svcs ...pick.Service[gin.HandlerFunc]) {
	gin2.OpenApi(engine, "", "")
	openApiLog()
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
			methodInfo := pick.GetMethodInfo(&method, preUrl, GinContextType)
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
			group.Handle(methodInfoExport.Method, methodInfoExport.Path[len(preUrl):], func(ctx *gin.Context) {
				ctxi := ginctx.FromRequest(ctx)
				defer ctxi.RootSpan().End()
				in2 := reflect.New(in2Type)
				err := binding.Bind(ctx, in2.Interface())
				if err != nil {
					ctx.JSON(http.StatusBadRequest, errcode.InvalidArgument.Msg(err.Error()))
					return
				}
				params := make([]reflect.Value, 3)
				params[0] = value
				if methodType.In(1).ConvertibleTo(GinContextType) {
					params[1] = reflect.ValueOf(ctxi)
				} else {
					params[1] = reflect.ValueOf(ctxi.Wrapper())
				}
				params[2] = in2
				result := methodValue.Call(params)
				pick.ResWriteReflect(ctx.Writer, ctxi.TraceID(), result)
			})
			methodInfo.Log()
			infos = append(infos, &apidoc2.ApiDocInfo{ApiInfo: methodInfoExport, Method: method.Type})
		}
		apidoc2.RegisterApiInfo(&apidoc2.GroupApiInfo{Describe: describe, Infos: infos})
	}
	pick.Registered()
}

func openApiLog() {
	pick.Log(http.MethodGet, apidoc.UriPrefix, "api文档列表")
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/swagger/*file", "swagger文档")
	pick.Log(http.MethodGet, apidoc.UriPrefix+"/markdown/*file", "markdown文档")
}

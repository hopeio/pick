/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/net/http/apidoc"
	basehigh "github.com/pb33f/libopenapi/datamodel/high/base"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"net/http"
	"reflect"
	"strings"
)

func Swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filePath, modName)
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			GenSwaggerApi(methodInfo.ApiInfo, doc, methodInfo.Method, groupApiInfo.Describe, methodInfo.Method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

func GenSwaggerApi(api *pick.ApiInfo, doc *v3high.Document, methodType reflect.Type, tag, dec string) {

	var pathItem *v3high.PathItem
	if doc.Paths.PathItems == nil {
		doc.Paths.PathItems = orderedmap.New[string, *v3high.PathItem]()
	}
	if path, ok := doc.Paths.PathItems.Get(api.Path); ok {
		pathItem = path
	} else {
		pathItem = new(v3high.PathItem)
	}

	//我觉得路径参数并没有那么值得非用不可
	parameters := make([]*v3high.Parameter, 0)
	numIn := methodType.NumIn()

	if numIn == 3 {
		if api.Method == http.MethodGet {
			InType := methodType.In(2).Elem()
			for j := 0; j < InType.NumField(); j++ {
				param := &v3high.Parameter{
					Name: InType.Field(j).Name,
					In:   "query",
				}
				parameters = append(parameters, param)
			}
		} else {
			reqName := methodType.In(2).Elem().Name()
			param := &v3high.Parameter{
				Name: reqName,
				In:   "body",
			}

			param.Schema = basehigh.CreateSchemaProxyRef(reqName)
			parameters = append(parameters, param)
			apidoc.DefinitionsApi(param.Schema.Schema(), reflect.New(methodType.In(2)).Elem().Interface())
		}
	}

	if !methodType.Out(0).Implements(pick.ErrorType) {
		var responses v3high.Responses
		responses.Codes = orderedmap.New[string, *v3high.Response]()
		var response v3high.Response
		response.Content = orderedmap.New[string, *v3high.MediaType]()
		responses.Default = &response
		response.Description = "一个成功的返回"
		schemaProxy := basehigh.CreateSchemaProxyRef(methodType.Out(0).Elem().Name())
		schema := schemaProxy.Schema()
		response.Content.Set("application/json", &v3high.MediaType{
			Schema: schemaProxy,
		})
		apidoc.DefinitionsApi(schema, reflect.New(methodType.Out(0)).Elem().Interface())
		responses.Codes.Set("200", &response)
		op := v3high.Operation{
			Summary:     api.Title,
			OperationId: api.Path + api.Method,
			Parameters:  parameters,
			Responses:   &responses,
		}

		var tags, desc []string
		tags = append(tags, tag, api.Createlog.Version)
		desc = append(desc, dec, api.Createlog.Log)
		for i := range api.Changelog {
			tags = append(tags, api.Changelog[i].Version)
			desc = append(desc, api.Changelog[i].Log)
		}
		op.Tags = tags
		op.Description = strings.Join(desc, "\n")

		switch api.Method {
		case http.MethodGet:
			pathItem.Get = &op
		case http.MethodPost:
			pathItem.Post = &op
		case http.MethodPut:
			pathItem.Put = &op
		case http.MethodDelete:
			pathItem.Delete = &op
		case http.MethodOptions:
			pathItem.Options = &op
		case http.MethodPatch:
			pathItem.Patch = &op
		case http.MethodHead:
			pathItem.Head = &op
		}
	}

	doc.Paths.PathItems.Set(api.Path, pathItem)
}

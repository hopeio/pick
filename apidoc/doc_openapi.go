/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/net/http/apidoc"
	"net/http"
	"reflect"
	"strings"
)

func Openapi(filePath, modName string) {
	doc := apidoc.GetDoc(filePath, modName)
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			GenOpenapi(methodInfo.ApiInfo, doc, methodInfo.Method, groupApiInfo.Describe, methodInfo.Method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

func GenOpenapi(api *pick.ApiInfo, doc *openapi3.T, methodType reflect.Type, tag, dec string) {
	for _, url := range api.Urls {
		var pathItem *openapi3.PathItem
		if doc.Paths != nil {
			if path := doc.Paths.Value(url.Path); path != nil {
				pathItem = path
			} else {
				pathItem = new(openapi3.PathItem)
			}
		} else {
			doc.Paths = openapi3.NewPaths()
			pathItem = new(openapi3.PathItem)
		}

		//我觉得路径参数并没有那么值得非用不可
		parameters := make([]*openapi3.ParameterRef, 0)
		numIn := methodType.NumIn()
		var requestBodyRef *openapi3.RequestBodyRef

		if numIn == 3 {
			if url.Method == http.MethodGet {
				InType := methodType.In(2).Elem()
				for j := 0; j < InType.NumField(); j++ {
					json := strings.Split(InType.Field(j).Tag.Get("json"), ",")[0]
					if json == "-" {
						continue
					}
					if json == "" {
						json = InType.Field(j).Name
					}
					param := &openapi3.ParameterRef{
						Value: &openapi3.Parameter{
							Name: json,
							In:   "query",
						},
					}
					parameters = append(parameters, param)
				}
			} else {
				reqName := methodType.In(2).Elem().Name()

				requestBody := &openapi3.RequestBody{Content: map[string]*openapi3.MediaType{"application/json": {Schema: openapi3.NewSchemaRef("#/components/schemas/"+reqName, nil)}}}
				requestBodyRef = &openapi3.RequestBodyRef{Value: requestBody}
				apidoc.AddComponent(reqName, reflect.New(methodType.In(2)).Elem().Interface())
			}
		}

		if !methodType.Out(0).Implements(pick.ErrorType) {
			responses := openapi3.NewResponses()
			response := responses.Default()
			mediaType := new(openapi3.MediaType)
			response.Value.Content = openapi3.NewContent()
			response.Value.Content["application/json"] = mediaType

			mediaType.Schema = openapi3.NewSchemaRef("#/components/schemas/"+methodType.Out(0).Elem().Name(), nil)
			apidoc.AddComponent(methodType.Out(0).Elem().Name(), reflect.New(methodType.Out(0)).Elem().Interface())
			title := api.Title
			if url.Remark != "" {
				title = title + "(" + url.Remark + ")"
			}
			op := openapi3.Operation{
				Summary:     title,
				OperationID: url.Method + url.Path,
				Parameters:  parameters,
				RequestBody: requestBodyRef,
				Responses:   responses,
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

			switch url.Method {
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

		doc.Paths.Set(url.Path, pathItem)
	}
}

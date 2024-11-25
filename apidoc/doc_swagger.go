/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"github.com/go-openapi/spec"
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/net/http/apidoc"
	reflecti "github.com/hopeio/utils/reflect"
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

func GenSwaggerApi(api *pick.ApiInfo, doc *spec.Swagger, methodType reflect.Type, tag, dec string) {
	if doc.Definitions == nil {
		doc.Definitions = make(map[string]spec.Schema)
	}

	var pathItem *spec.PathItem
	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[api.Path]; ok {
			pathItem = &path
		} else {
			pathItem = new(spec.PathItem)
		}
	} else {
		doc.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = new(spec.PathItem)
	}

	//我觉得路径参数并没有那么值得非用不可
	parameters := make([]spec.Parameter, 0)
	numIn := methodType.NumIn()

	if numIn == 3 {
		if api.Method == http.MethodGet {
			InType := methodType.In(2).Elem()
			for j := 0; j < InType.NumField(); j++ {
				param := spec.Parameter{
					ParamProps: spec.ParamProps{
						Name: InType.Field(j).Name,
						In:   "query",
					},
				}
				parameters = append(parameters, param)
			}
		} else {
			reqName := methodType.In(2).Elem().Name()
			param := spec.Parameter{
				ParamProps: spec.ParamProps{
					Name: reqName,
					In:   "body",
				},
			}

			param.Schema = new(spec.Schema)
			param.Schema.Ref = spec.MustCreateRef("#/definitions/" + reqName)
			parameters = append(parameters, param)
			DefinitionsApi(doc.Definitions, reflect.New(methodType.In(2)).Elem().Interface(), nil)
		}
	}

	if !methodType.Out(0).Implements(pick.ErrorType) {
		var responses spec.Responses
		responses.StatusCodeResponses = make(map[int]spec.Response)
		response := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
		response.Schema.Ref = spec.MustCreateRef("#/definitions/" + methodType.Out(0).Elem().Name())
		response.Description = "一个成功的返回"
		DefinitionsApi(doc.Definitions, reflect.New(methodType.Out(0)).Elem().Interface(), nil)
		responses.StatusCodeResponses[200] = response
		op := spec.Operation{
			OperationProps: spec.OperationProps{
				Summary:    api.Title,
				ID:         api.Path + api.Method,
				Parameters: parameters,
				Responses:  &responses,
			},
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

	doc.Paths.Paths[api.Path] = *pathItem
}

func DefinitionsApi(definitions map[string]spec.Schema, v interface{}, exclude []string) {
	schema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       []string{"object"},
			Properties: make(map[string]spec.Schema),
		},
	}

	body := reflect.TypeOf(v).Elem()
	var typ, subFieldName string
	var arraySubType string
	for i := 0; i < body.NumField(); i++ {
		json := strings.Split(body.Field(i).Tag.Get("json"), ",")[0]
		if json == "" || json == "-" {
			continue
		}
		fieldType := body.Field(i).Type
		switch fieldType.Kind() {
		case reflect.Struct:
			typ = "object"
			v = reflect.New(fieldType).Interface()
			subFieldName = fieldType.Name()
		case reflect.Ptr:
			typ = "object"
			v = reflect.New(fieldType.Elem()).Interface()
			subFieldName = fieldType.Elem().Name()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Array, reflect.Slice:
			typ = "array"
			subType := reflecti.DerefType(fieldType)
			subFieldName = subType.Name()
			switch subType.Kind() {
			case reflect.Struct, reflect.Ptr, reflect.Array, reflect.Slice:
				v = reflect.New(subType).Interface()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				arraySubType = "integer"
			case reflect.Float32, reflect.Float64:
				arraySubType = "number"
			case reflect.String:
				arraySubType = "string"
			case reflect.Bool:
				arraySubType = "boolean"
			}
		case reflect.Float32, reflect.Float64:
			typ = "number"
		case reflect.String:
			typ = "string"
		case reflect.Bool:
			typ = "boolean"
		}
		subSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{typ},
			},
		}
		if typ == "object" {
			subSchema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		if typ == "array" {
			subSchema.Items = new(spec.SchemaOrArray)
			subSchema.Items.Schema = &spec.Schema{}
			if arraySubType == "" {
				subSchema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
				DefinitionsApi(definitions, v, nil)
			} else {
				subSchema.Items.Schema.Type = []string{arraySubType}
			}

		}
		schema.Properties[json] = subSchema
	}
	definitions[body.Name()] = schema
}

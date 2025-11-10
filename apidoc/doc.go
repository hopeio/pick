/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/apidoc"
	"github.com/hopeio/gox/reflect/structtag"
	"github.com/hopeio/pick"
)

func DocList(w http.ResponseWriter, r *http.Request) {
	modName := r.URL.Query().Get("modName")
	if modName == "" {
		modName = "api"
	}
	Openapi(apidoc.Dir, modName)
	apidoc.DocList(w, r)
}

type ApiDocInfo struct {
	ApiInfo *pick.ApiInfo
	Method  reflect.Type
}

type GroupApiInfo struct {
	Describe string
	Infos    []*ApiDocInfo
}

var groupApiInfos []*GroupApiInfo

func RegisterApiInfo(apiInfo *GroupApiInfo) {
	groupApiInfos = append(groupApiInfos, apiInfo)
}

var Doc *openapi3.T

// 参数为路径和格式
func GetDoc(realPath, modName string) *openapi3.T {
	if Doc != nil {
		return Doc
	}
	api := NewAPI(modName)
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			GenOpenapi(methodInfo, api, groupApiInfo.Describe)
		}
	}
	spec, err := api.Spec()
	if err != nil {
		log.Error(err)
	}
	Doc = spec
	return Doc
}

func WriteToFile(realPath, modName string) {
	if Doc == nil {
		return
	}
	if realPath == "" {
		realPath = "."
	}

	realPath = realPath + modName
	err := os.MkdirAll(realPath, os.ModePerm)
	if err != nil {
		log.Error(err)
	}

	realPath = filepath.Join(realPath, modName+apidoc.OpenapiEXT)

	if _, err := os.Stat(realPath); err == nil {
		os.Remove(realPath)
	}
	var file *os.File
	file, err = os.Create(realPath)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(Doc)
	if err != nil {
		log.Error(err)
	}

	/*b, err := yaml.Marshal(swag.ToDynamicJSON(Doc))
	  if err != nil {
	  	log.Error(err)
	  }
	  if _, err := file.Write(b); err != nil {
	  	log.Error(err)
	  }*/

	Doc = nil
}

func Openapi(filePath, modName string) {
	GetDoc(filePath, modName)
	WriteToFile(filePath, modName)
}

func GenOpenapi(methodInfo *ApiDocInfo, api *API, dec string) {
	for _, route := range methodInfo.ApiInfo.Routes {
		r := api.Route(route.Method, route.Path)
		numIn := methodInfo.Method.NumIn()
		if numIn == 3 {
			InType := methodInfo.Method.In(2).Elem()
			for j := 0; j < InType.NumField(); j++ {
				tags, err := structtag.Parse(string(InType.Field(j).Tag))
				if err != nil {
					return
				}
				if uri := tags.MustGet("uri"); uri.Value != "" && uri.Value != "-" {
					r.HasPathParameter(uri.Value, PathParam{
						Description:       tags.MustGet("comment").Value,
						Regexp:            "",
						Type:              Type(InType.Field(j).Type),
						ApplyCustomSchema: nil,
					})
				}
				if uri := tags.MustGet("path"); uri.Value != "" && uri.Value != "-" {
					r.HasPathParameter(uri.Value, PathParam{
						Description:       tags.MustGet("comment").Value,
						Regexp:            "",
						Type:              Type(InType.Field(j).Type),
						ApplyCustomSchema: nil,
					})
				}
				if query := tags.MustGet("query"); query.Value != "" && query.Value != "-" {
					r.HasQueryParameter(query.Value, QueryParam{
						Description:       tags.MustGet("comment").Value,
						Regexp:            "",
						Type:              Type(InType.Field(j).Type),
						ApplyCustomSchema: nil,
					})
				}
				if uri := tags.MustGet("json"); uri.Value != "" && uri.Value != "-" {
					r.HasRequestModel(Model{Type: InType})
				}
			}
		}
		r.HasTags([]string{dec, methodInfo.ApiInfo.Title + route.Remark})
		if methodInfo.Method.Out(0).Implements(pick.ErrorType) {
			r.HasResponseModel(http.StatusOK, Model{Type: pick.ErrRepType})
		} else {
			r.HasResponseModel(http.StatusOK, Model{Type: methodInfo.Method.Out(0)})
		}
		r.HasResponseModel(http.StatusBadRequest, Model{Type: pick.ErrRepType})

	}
}

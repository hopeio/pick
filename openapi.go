/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/openapi"
)

func DocList(w http.ResponseWriter, r *http.Request) {
	modName := r.URL.Query().Get("modName")
	if modName == "" {
		modName = "api"
	}

	openapi.DocList(w, r)
}

type ApiDocInfo struct {
	ApiInfo *ApiInfo
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

func Openapi(docDir, modName string) {
	doc := GetDoc(modName)
	openapi.WriteToFile(docDir, modName, doc)
}

var Doc *openapi3.T

// 参数为路径和格式
func GetDoc(modName string) *openapi3.T {
	if Doc != nil {
		return Doc
	}
	api := openapi.NewAPI(modName)
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

func GenOpenapi(methodInfo *ApiDocInfo, api *openapi.API, dec string) {
	for _, route := range methodInfo.ApiInfo.Routes {
		r := api.Route(route.Method, route.Path)
		numIn := methodInfo.Method.NumIn()
		if numIn == 3 {
			r.HasRequest(openapi.Model{Type: methodInfo.Method.In(2).Elem()})
		}
		r.HasTags([]string{dec, methodInfo.ApiInfo.Title + route.Remark})
		if methodInfo.Method.Out(0).Implements(ErrorType) {
			r.HasResponseModel(http.StatusOK, openapi.Model{Type: ErrRespType})
		} else {
			r.HasResponseModel(http.StatusOK, openapi.Model{Type: methodInfo.Method.Out(0)})
		}
		r.HasResponseModel(http.StatusBadRequest, openapi.Model{Type: ErrRespType})

	}
}

/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/apidoc"
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

func Openapi(filePath, modName string) {
	doc := GetDoc(modName)
	apidoc.WriteToFile(filePath, modName, doc)
}

var Doc *openapi3.T

// 参数为路径和格式
func GetDoc(modName string) *openapi3.T {
	if Doc != nil {
		return Doc
	}
	api := apidoc.NewAPI(modName)
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

func GenOpenapi(methodInfo *ApiDocInfo, api *apidoc.API, dec string) {
	for _, route := range methodInfo.ApiInfo.Routes {
		r := api.Route(route.Method, route.Path)
		numIn := methodInfo.Method.NumIn()
		if numIn == 3 {
			r.HasRequest(apidoc.Model{Type: methodInfo.Method.In(2).Elem()})
		}
		r.HasTags([]string{dec, methodInfo.ApiInfo.Title + route.Remark})
		if methodInfo.Method.Out(0).Implements(pick.ErrorType) {
			r.HasResponseModel(http.StatusOK, apidoc.Model{Type: pick.ErrRespType})
		} else {
			r.HasResponseModel(http.StatusOK, apidoc.Model{Type: methodInfo.Method.Out(0)})
		}
		r.HasResponseModel(http.StatusBadRequest, apidoc.Model{Type: pick.ErrRespType})

	}
}

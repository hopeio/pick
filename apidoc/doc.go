/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"github.com/hopeio/pick"
	"github.com/hopeio/utils/net/http/apidoc"
	"net/http"
	"reflect"
)

func DocList(w http.ResponseWriter, r *http.Request) {
	modName := r.URL.Query().Get("modName")
	if modName == "" {
		modName = "api"
	}
	Markdown(apidoc.Dir, modName)
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

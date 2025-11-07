/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/http/apidoc"
	"github.com/hopeio/pick"
	"gopkg.in/yaml.v3"
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
	if realPath == "" {
		realPath = "."
	}

	realPath = realPath + modName
	err := os.MkdirAll(realPath, os.ModePerm)
	if err != nil {
		log.Error(err)
	}

	realPath = filepath.Join(realPath, modName+apidoc.OpenapiEXT)

	apiType := filepath.Ext(realPath)

	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		Doc = newSpec(modName)
		return Doc
	} else {
		file, err := os.Open(realPath)
		if err != nil {
			log.Error(err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			log.Error(err)
		}
		/*var buf bytes.Buffer
		err = json.Compact(&buf, data)
		if err != nil {
			ulog.Error(err)
		}*/
		if apiType == ".json" {
			err = json.Unmarshal(data, &Doc)
			if err != nil {
				log.Error(err)
			}
		} else {
			//var v map[string]interface{}//子类型 json: unsupported type: map[interface{}]interface{}
			//var v interface{} //json: unsupported type: map[interface{}]interface{}
			err = yaml.Unmarshal(data, &Doc)
			if err != nil {
				log.Error(err)
			}
		}
	}
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

func NilDoc() {
	Doc = nil
}

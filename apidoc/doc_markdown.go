/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package apidoc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	reflectx "github.com/hopeio/gox/reflect"
	"github.com/hopeio/gox/reflect/mock"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/hopeio/gox/validation/validator"
)

// 有swagger,有没有必要做
func Markdown(filePath, modName string) {
	buf, err := genFile(filePath, modName)
	if err != nil {
		log.Println(err)
	}
	defer buf.Close()
	fmt.Fprintln(buf, "[TOC]")
	if modName != "" {
		fmt.Fprintf(buf, "# %s接口文档  \n", modName)
		fmt.Fprintln(buf, "----------")
	}
	for _, groupApiInfo := range groupApiInfos {
		fmt.Fprintf(buf, "# %s  \n", groupApiInfo.Describe)
		fmt.Fprintln(buf, "----------")
		for _, methodInfo := range groupApiInfo.Infos {
			//title
			apiInfo := methodInfo.ApiInfo
			for _, url := range apiInfo.Routes {
				title := apiInfo.Title
				if url.Remark != "" {
					title = title + "(" + url.Remark + ")"
				}
				if apiInfo.Deprecated != nil {
					fmt.Fprintf(buf, "## ~~%s-%s(废弃)(`%s`)~~  \n", title, apiInfo.Deprecated.Version, url.Path)
				} else {
					log := apiInfo.Createlog
					if len(apiInfo.Changelog) > 0 {
						log = apiInfo.Changelog[len(apiInfo.Changelog)-1]
					}
					fmt.Fprintf(buf, "## %s-%s(`%s`)  \n", title, log.Version, url.Path)
				}
				//api
				fmt.Fprintf(buf, "**%s** `%s` _(Principal %s)_  \n", url.Method, url.Path, apiInfo.GetPrincipal())

				fmt.Fprint(buf, "### 接口记录  \n")
				fmt.Fprint(buf, "|版本|操作|时间|负责人|日志|  \n")
				fmt.Fprint(buf, "| :----: | :----: | :----: | :----: | :----: |  \n")
				fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", apiInfo.Createlog.Version, "创建", apiInfo.Createlog.Date, apiInfo.Createlog.Auth, apiInfo.Createlog.Desc)
				if len(apiInfo.Changelog) != 0 || apiInfo.Deprecated != nil {
					for _, clog := range apiInfo.Changelog {
						fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", clog.Version, "变更", clog.Date, clog.Auth, clog.Desc)
					}
					if apiInfo.Deprecated != nil {
						fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", apiInfo.Deprecated.Version, "删除", apiInfo.Deprecated.Date, apiInfo.Deprecated.Auth, apiInfo.Deprecated.Desc)
					}
				}

				fmt.Fprint(buf, "### 参数信息  \n")
				if methodInfo.Method.NumIn() == 3 {
					fmt.Fprint(buf, "|字段名称|字段类型|字段描述|校验要求|  \n")
					fmt.Fprint(buf, "| :----  | :----: | :----: | :----: |  \n")
					params := getParamTable(methodInfo.Method.In(2).Elem(), "")
					for i := range params {
						fmt.Fprintf(buf, "|%s|%s|%s|%s|  \n", params[i].json, params[i].typ, params[i].comment, params[i].validator)
					}

				} else {
					fmt.Fprint(buf, "无需参数")
				}
				fmt.Fprint(buf, "__请求示例__  \n")
				fmt.Fprint(buf, "```json  \n")
				newParam := reflect.New(methodInfo.Method.In(2).Elem()).Interface()
				mock.Mock(newParam)
				data, _ := json.MarshalIndent(newParam, "", "\t")
				fmt.Fprint(buf, string(data), "  \n")
				fmt.Fprint(buf, "```  \n")
				fmt.Fprint(buf, "### 返回信息  \n")
				fmt.Fprint(buf, "|字段名称|字段类型|字段描述|  \n")
				fmt.Fprint(buf, "| :----  | :----: | :----: | \n")
				params := getParamTable(methodInfo.Method.Out(0).Elem(), "")
				for i := range params {
					fmt.Fprintf(buf, "|%s|%s|%s|  \n", params[i].json, params[i].typ, params[i].comment)
				}
				fmt.Fprint(buf, "__返回示例__  \n")
				fmt.Fprint(buf, "```json  \n")
				newRes := reflect.New(methodInfo.Method.Out(0).Elem()).Interface()
				mock.Mock(newRes)
				data, _ = json.MarshalIndent(newRes, "", "\t")
				fmt.Fprint(buf, string(data), "  \n")
				fmt.Fprint(buf, "```  \n")
			}
		}
	}
}

func genFile(filePath, modName string) (*os.File, error) {

	filePath = filePath + modName

	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	filePath = filepath.Join(filePath, modName+".apidoc.md")

	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}
	var file *os.File
	file, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type ParamTable struct {
	json, comment, typ, validator string
}

func getParamTable(param reflect.Type, pre string) []*ParamTable {
	param = reflectx.DerefType(param)
	newParam := reflect.New(param).Interface()
	var res []*ParamTable
	for i := range param.NumField() {
		/*		if param.AssignableTo(reflect.TypeOf(response.File{})) {
				return "下载文件"
			}*/
		var p ParamTable
		field := param.Field(i)
		if field.Anonymous {
			continue
		}
		json := strings.Split(field.Tag.Get("json"), ",")[0]
		if json == "-" {
			continue
		}
		if json == "" {
			p.json = pre + stringsx.SnakeToCamel(json)
		} else {
			p.json = pre + json
		}
		p.comment = field.Tag.Get("comment")
		if p.comment == "-" {
			p.comment = p.json
		}
		p.typ = getJsType(field.Type)
		if valid := validator.TransError(validator.Validator.StructPartial(newParam, field.Name)); valid != "" {
			p.validator = valid[len(p.comment):]
		}
		if p.typ == "object" || p.typ == "[]object" {
			p.json = "**" + p.json + "**"
			res = append(res, &p)
			sub := getParamTable(field.Type, json+".")
			res = append(res, sub...)
		} else {
			res = append(res, &p)
		}
	}
	return res
}

func getJsType(typ reflect.Type) string {
	t := time.Time{}
	if typ == reflect.TypeOf(t) || typ == reflect.TypeOf(&t) {
		return "date"
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Array, reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return "string"
		}
		return "[]" + getJsType(typ.Elem())
	case reflect.Ptr:
		return getJsType(typ.Elem())
	case reflect.Struct:
		return "object"
	case reflect.Bool:
		return "boolean"
	}
	return "string"
}

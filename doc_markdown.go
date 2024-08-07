package pick

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	reflecti "github.com/hopeio/utils/reflect"
	"github.com/hopeio/utils/reflect/mock"
	stringsi "github.com/hopeio/utils/strings"
	"github.com/hopeio/utils/validation/validator"
)

type ApiDocInfo struct {
	ApiInfo *apiInfo
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
			apiinfo := methodInfo.ApiInfo
			if apiinfo.deprecated != nil {
				fmt.Fprintf(buf, "## ~~%s-v%d(废弃)(`%s`)~~  \n", apiinfo.title, apiinfo.version, apiinfo.path)
			} else {
				fmt.Fprintf(buf, "## %s-v%d(`%s`)  \n", apiinfo.title, apiinfo.version, apiinfo.path)
			}
			//api
			fmt.Fprintf(buf, "**%s** `%s` _(Principal %s)_  \n", apiinfo.method, apiinfo.path, apiinfo.getPrincipal())

			fmt.Fprint(buf, "### 接口记录  \n")
			fmt.Fprint(buf, "|版本|操作|时间|负责人|日志|  \n")
			fmt.Fprint(buf, "| :----: | :----: | :----: | :----: | :----: |  \n")
			fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", apiinfo.createlog.version, "创建", apiinfo.createlog.date, apiinfo.createlog.auth, apiinfo.createlog.log)
			if len(apiinfo.changelog) != 0 || apiinfo.deprecated != nil {
				for _, clog := range apiinfo.changelog {
					fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", clog.version, "变更", clog.date, clog.auth, clog.log)
				}
				if apiinfo.deprecated != nil {
					fmt.Fprintf(buf, "|%s|%s|%s|%s|%s|  \n", apiinfo.deprecated.version, "删除", apiinfo.deprecated.date, apiinfo.deprecated.auth, apiinfo.deprecated.log)
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
	param = reflecti.DerefType(param)
	newParam := reflect.New(param).Interface()
	var res []*ParamTable
	for i := 0; i < param.NumField(); i++ {
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
			p.json = pre + stringsi.SnakeToCamel(json)
		} else {
			p.json = pre + json
		}
		p.comment = field.Tag.Get("comment")
		if p.comment == "-" {
			p.comment = p.json
		}
		p.typ = getJsType(field.Type)
		if valid := validator.Trans(validator.Validator.StructPartial(newParam, field.Name)); valid != "" {
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

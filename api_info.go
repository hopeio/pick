/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"context"
	"errors"
	"github.com/hopeio/utils/log"
	"net/http"
	"reflect"
	"strings"
	"unsafe"
)

const Template = `
func (*UserService) Add(ctx *model.Context, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("2.0.1", "jyb", "2019/12/16", "修改测试").End()
	})

	return &response.TinyRep{Msg: req.Name}, nil
}
`

var (
	ContextType  = reflect.TypeOf((*context.Context)(nil)).Elem()
	ContextValue = reflect.ValueOf(context.Background())
	ErrorType    = reflect.TypeOf((*error)(nil)).Elem()
)

type url struct {
	Path, Method, Remark string
}

func Get(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodGet}}}
}
func Post(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodPost}}}
}
func Put(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodPut}}}
}
func Delete(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodDelete}}}
}
func Patch(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodPatch}}}
}
func Trace(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodTrace}}}
}
func Head(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodHead}}}
}
func Options(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodOptions}}}
}
func Connect(path string) *apiInfo {
	return &apiInfo{urls: []url{{Path: path, Method: http.MethodConnect}}}
}

type apiInfo struct {
	urls       []url
	title      string
	version    int
	changelog  []changelog
	createlog  changelog
	deprecated *changelog
}

type changelog struct {
	Version, Auth, Date, Log string
}

func (api *apiInfo) Get(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodGet})
	return api
}

func (api *apiInfo) GetRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodGet, Remark: remark})
	return api
}

func (api *apiInfo) Post(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPost})
	return api
}

func (api *apiInfo) PostRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPost, Remark: remark})
	return api
}

func (api *apiInfo) Put(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPut})
	return api
}

func (api *apiInfo) PutRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPut, Remark: remark})
	return api
}

func (api *apiInfo) Delete(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodDelete})
	return api
}

func (api *apiInfo) DeleteRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodDelete, Remark: remark})
	return api
}

func (api *apiInfo) Patch(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPatch})
	return api
}

func (api *apiInfo) PatchRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodPatch, Remark: remark})
	return api
}

func (api *apiInfo) Trace(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodTrace})
	return api
}

func (api *apiInfo) TraceRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodTrace, Remark: remark})
	return api
}

func (api *apiInfo) Head(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodHead})
	return api
}
func (api *apiInfo) HeadRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodHead, Remark: remark})
	return api
}
func (api *apiInfo) Options(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodOptions})
	return api
}

func (api *apiInfo) OptionsRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodOptions, Remark: remark})
	return api
}
func (api *apiInfo) Connect(path string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodConnect})
	return api
}

func (api *apiInfo) ConnectRemark(path, remark string) *apiInfo {
	api.urls = append(api.urls, url{Path: path, Method: http.MethodConnect, Remark: remark})
	return api
}

func (api *apiInfo) ChangeLog(v, auth, date, log string) *apiInfo {
	v = version(v)
	api.changelog = append(api.changelog, changelog{v, auth, date, log})
	return api
}

func version(v string) string {
	if v[0] != 'v' {
		return "v" + v
	}
	return v
}

func (api *apiInfo) CreateLog(v, auth, date, log string) *apiInfo {
	if api.createlog.Version != "" {
		panic("创建记录只允许一条")
	}
	v = version(v)
	api.createlog = changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Title(d string) *apiInfo {
	api.title = d
	return api
}

func (api *apiInfo) Version(v int) *apiInfo {
	api.version = v
	return api
}

func (api *apiInfo) Deprecated(v, auth, date, log string) *apiInfo {
	v = version(v)
	api.deprecated = &changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) End() {
	panic(api)
}

func (api *apiInfo) Check() error {
	if len(api.urls) == 0 || api.urls[0].Path == "" || api.urls[0].Method == "" || api.title == "" {
		return errors.New("接口路径,方法,描述均为必填")
	}
	return nil
}

func (api *apiInfo) Log() {
	for _, url := range api.urls {
		Log(url.Method, url.Path, api.title)
	}
}

func (api *apiInfo) Export() *ApiInfo {
	return (*ApiInfo)(unsafe.Pointer(api))
}

type ApiInfo struct {
	Urls       []url
	Title      string
	Version    int
	Changelog  []changelog
	Createlog  changelog
	Deprecated *changelog
}

// 获取负责人
func (api *ApiInfo) GetPrincipal() string {
	if len(api.Changelog) == 0 {
		return api.Createlog.Auth
	}
	if api.Deprecated != nil {
		return api.Deprecated.Auth
	}
	return api.Changelog[len(api.Changelog)-1].Auth
}

// recover捕捉panic info
func GetMethodInfo(method *reflect.Method, preUrl string, httpContext reflect.Type) (info *apiInfo) {
	if method.Name == "Service" {
		return nil
	}
	if !strings.HasPrefix(method.Name, prefix) {
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				//_,_, info.Version = ParseMethodName(Method.Name)
				if v.version == 0 {
					v.version = 1
				}
				for i := range v.urls {
					v.urls[i].Path = preUrl + v.urls[i].Path
				}

				info = v
			} else {
				log.Error(err)
			}
		}
	}()
	methodValue := method.Func
	methodType := methodValue.Type()
	numIn := methodType.NumIn()
	numOut := methodType.NumOut()
	var err error
	defer func() {
		if err != nil {
			log.Debugf("未注册: %s 原因:%v", method.Name, err)
		}
	}()

	if numIn != 3 {
		//err = errors.New("method参数必须为两个")
		return
	}
	if numOut != 2 {
		//err = errors.New("method返回值必须为两个")
		return
	}

	if !methodType.In(1).ConvertibleTo(httpContext) && !methodType.In(1).Implements(ContextType) {
		err = errors.New("service第一个参数必须为*httpctx.Context类型或context.Context")
		return
	}
	if !methodType.Out(1).Implements(ErrorType) {
		err = errors.New("service第二个返回值必须为error类型")
		return
	}
	params := make([]reflect.Value, numIn)
	params[0] = reflect.New(methodType.In(0).Elem())
	if methodType.In(1).ConvertibleTo(httpContext) {
		params[1] = reflect.New(methodType.In(1).Elem())
	} else {
		params[1] = ContextValue
	}
	params[2] = reflect.New(methodType.In(2).Elem())
	methodValue.Call(params)
	return nil
}

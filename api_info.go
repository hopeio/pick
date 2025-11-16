/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/hopeio/gox/log"
	unsafei "github.com/hopeio/gox/unsafe"
)

const Template = `
func (*UserService) Add(ctx *httpreq.Context, req *model.SignupReq) (*response.TinyResp, error) {
	pick.Api(func() {
		pick.Post("").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("2.0.1", "jyb", "2019/12/16", "修改测试").End()
	})

	return &response.TinyResp{Msg: req.Name}, nil
}
`

var (
	ContextType  = reflect.TypeOf((*context.Context)(nil)).Elem()
	ContextValue = reflect.ValueOf(context.Background())
	ErrorType    = reflect.TypeOf((*error)(nil)).Elem()
)

type Route struct {
	Path, Method, Remark string
}

func Get(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodGet}}}
}
func Post(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodPost}}}
}
func Put(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodPut}}}
}
func Delete(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodDelete}}}
}
func Patch(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodPatch}}}
}
func Trace(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodTrace}}}
}
func Head(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodHead}}}
}
func Options(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodOptions}}}
}
func Connect(path string) *apiInfo {
	return &apiInfo{routes: []Route{{Path: path, Method: http.MethodConnect}}}
}

type apiInfo struct {
	routes []Route
	title  string
}

type Changelog struct {
	Version, Auth, Date, Desc string
}

func (api *apiInfo) Get(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodGet})
	return api
}

func (api *apiInfo) GetRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodGet, Remark: remark})
	return api
}

func (api *apiInfo) Post(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPost})
	return api
}

func (api *apiInfo) PostRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPost, Remark: remark})
	return api
}

func (api *apiInfo) Put(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPut})
	return api
}

func (api *apiInfo) PutRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPut, Remark: remark})
	return api
}

func (api *apiInfo) Delete(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodDelete})
	return api
}

func (api *apiInfo) DeleteRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodDelete, Remark: remark})
	return api
}

func (api *apiInfo) Patch(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPatch})
	return api
}

func (api *apiInfo) PatchRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPatch, Remark: remark})
	return api
}

func (api *apiInfo) Trace(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodTrace})
	return api
}

func (api *apiInfo) TraceRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodTrace, Remark: remark})
	return api
}

func (api *apiInfo) Head(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodHead})
	return api
}
func (api *apiInfo) HeadRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodHead, Remark: remark})
	return api
}
func (api *apiInfo) Options(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodOptions})
	return api
}

func (api *apiInfo) OptionsRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodOptions, Remark: remark})
	return api
}
func (api *apiInfo) Connect(path string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodConnect})
	return api
}

func (api *apiInfo) ConnectRemark(path, remark string) *apiInfo {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodConnect, Remark: remark})
	return api
}

func version(v string) string {
	if v[0] != 'v' {
		return "v" + v
	}
	return v
}

func (api *apiInfo) Title(d string) *apiInfo {
	api.title = d
	return api
}

func (api *apiInfo) End() {
	panic(api)
}

func (api *apiInfo) Check() error {
	if len(api.routes) == 0 || api.routes[0].Path == "" || api.routes[0].Method == "" || api.title == "" {
		return errors.New("path,method,title is required")
	}
	return nil
}

func (api *apiInfo) Log() {
	for _, url := range api.routes {
		Log(url.Method, url.Path, api.title)
	}
}

func (api *apiInfo) Export() *ApiInfo {
	return (*ApiInfo)(unsafe.Pointer(api))
}

type ApiInfo struct {
	Routes []Route
	Title  string
}

// recover捕捉panic info
func GetMethodInfo[T any](method *reflect.Method, preUrl string, httpContext reflect.Type) (info *apiInfo) {
	if method.Name == "Service" {
		return nil
	}
	if prefix != "" && !strings.HasPrefix(method.Name, prefix) {
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				for i := range v.routes {
					v.routes[i].Path = preUrl + v.routes[i].Path
				}
				info = v
			} else if v, ok := err.(*apiInfo); ok {
				for i := range v.routes {
					v.routes[i].Path = preUrl + v.routes[i].Path
				}
				info = unsafei.Cast[apiInfo](v)
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
			log.Debugf("unregisted: %s reason:%v", method.Name, err)
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
		err = errors.New("service first argument should be *httpctx.Context type or context.Context")
		return
	}
	if !methodType.Out(1).Implements(ErrorType) && methodType.Out(1) != ErrRespType {
		err = errors.New("service second return should be error type")
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

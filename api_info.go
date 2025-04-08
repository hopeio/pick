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
	unsafei "github.com/hopeio/utils/unsafe"
	"net/http"
	"reflect"
	"strings"
	"unsafe"
)

const Template = `
func (*UserService) Add(ctx *httpreq.Context, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("").
			Title("用户注册").
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

type Route struct {
	Path, Method, Remark string
}

type none = struct{}

func Middleware[T any](middleware ...T) *apiInfo[T] {
	return &apiInfo[T]{middlewares: middleware}
}

func Get(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodGet}}}
}
func Post(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodPost}}}
}
func Put(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodPut}}}
}
func Delete(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodDelete}}}
}
func Patch(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodPatch}}}
}
func Trace(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodTrace}}}
}
func Head(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodHead}}}
}
func Options(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodOptions}}}
}
func Connect(path string) *apiInfo[none] {
	return &apiInfo[none]{routes: []Route{{Path: path, Method: http.MethodConnect}}}
}

type apiInfo[T any] struct {
	routes      []Route
	middlewares []T
	title       string
	createlog   Changelog
	changelog   []Changelog
	deprecated  *Changelog
}

type Changelog struct {
	Version, Auth, Date, Desc string
}

func (api *apiInfo[T]) Get(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodGet})
	return api
}

func (api *apiInfo[T]) GetRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodGet, Remark: remark})
	return api
}

func (api *apiInfo[T]) Post(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPost})
	return api
}

func (api *apiInfo[T]) PostRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPost, Remark: remark})
	return api
}

func (api *apiInfo[T]) Put(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPut})
	return api
}

func (api *apiInfo[T]) PutRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPut, Remark: remark})
	return api
}

func (api *apiInfo[T]) Delete(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodDelete})
	return api
}

func (api *apiInfo[T]) DeleteRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodDelete, Remark: remark})
	return api
}

func (api *apiInfo[T]) Patch(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPatch})
	return api
}

func (api *apiInfo[T]) PatchRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodPatch, Remark: remark})
	return api
}

func (api *apiInfo[T]) Trace(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodTrace})
	return api
}

func (api *apiInfo[T]) TraceRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodTrace, Remark: remark})
	return api
}

func (api *apiInfo[T]) Head(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodHead})
	return api
}
func (api *apiInfo[T]) HeadRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodHead, Remark: remark})
	return api
}
func (api *apiInfo[T]) Options(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodOptions})
	return api
}

func (api *apiInfo[T]) OptionsRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodOptions, Remark: remark})
	return api
}
func (api *apiInfo[T]) Connect(path string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodConnect})
	return api
}

func (api *apiInfo[T]) ConnectRemark(path, remark string) *apiInfo[T] {
	api.routes = append(api.routes, Route{Path: path, Method: http.MethodConnect, Remark: remark})
	return api
}

func (api *apiInfo[T]) ChangeLog(v, auth, date, log string) *apiInfo[T] {
	v = version(v)
	api.changelog = append(api.changelog, Changelog{v, auth, date, log})
	return api
}

func version(v string) string {
	if v[0] != 'v' {
		return "v" + v
	}
	return v
}

func (api *apiInfo[T]) CreateLog(v, auth, date, log string) *apiInfo[T] {
	if api.createlog.Version != "" {
		panic("only one createlog is allowed")
	}
	v = version(v)
	api.createlog = Changelog{v, auth, date, log}
	return api
}

func (api *apiInfo[T]) Title(d string) *apiInfo[T] {
	api.title = d
	return api
}

func (api *apiInfo[T]) Deprecated(v, auth, date, log string) *apiInfo[T] {
	v = version(v)
	api.deprecated = &Changelog{v, auth, date, log}
	return api
}

func (api *apiInfo[T]) End() {
	panic(api)
}

func (api *apiInfo[T]) Check() error {
	if len(api.routes) == 0 || api.routes[0].Path == "" || api.routes[0].Method == "" || api.title == "" {
		return errors.New("path,method,title is required")
	}
	return nil
}

func (api *apiInfo[T]) Log() {
	for _, url := range api.routes {
		Log(url.Method, url.Path, api.title)
	}
}

func (api *apiInfo[T]) Export() *ApiInfo {
	return (*ApiInfo)(unsafe.Pointer(api))
}

type ApiInfo struct {
	Routes      []Route
	Middlewares []none
	Title       string
	Createlog   Changelog
	Changelog   []Changelog
	Deprecated  *Changelog
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
func GetMethodInfo[T any](method *reflect.Method, preUrl string, httpContext reflect.Type) (info *apiInfo[T]) {
	if method.Name == "Service" {
		return nil
	}
	if prefix != "" && !strings.HasPrefix(method.Name, prefix) {
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo[T]); ok {
				for i := range v.routes {
					v.routes[i].Path = preUrl + v.routes[i].Path
				}
				info = v
			} else if v, ok := err.(*apiInfo[none]); ok {
				for i := range v.routes {
					v.routes[i].Path = preUrl + v.routes[i].Path
				}
				info = unsafei.Cast[apiInfo[T]](v)
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
	if !methodType.Out(1).Implements(ErrorType) && methodType.Out(1) != ErrRepType {
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

package pick

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hopeio/utils/net/http/apidoc"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/hopeio/utils/log"
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

type apiInfo struct {
	path, method, title string
	version             int
	changelog           []changelog
	createlog           changelog
	deprecated          *changelog
	middleware          []http.HandlerFunc
}

type changelog struct {
	version, auth, date, log string
}

func Get(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodGet}
}
func Post(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPost}
}
func Put(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPut}
}
func Delete(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodDelete}
}
func Patch(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPatch}
}
func Trace(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodTrace}
}
func Head(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodHead}
}
func Options(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodOptions}
}
func Connect(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodConnect}
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
	if api.createlog.version != "" {
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

func (api *apiInfo) Middleware(m ...http.HandlerFunc) *apiInfo {
	api.middleware = m
	return api
}

func (api *apiInfo) End() {
	panic(api)
}

func (api *apiInfo) Check() error {
	if api.path == "" || api.method == "" || api.title == "" || api.createlog.version == "" {
		return errors.New("接口路径,方法,描述,创建日志均为必填")
	}
	return nil
}

func (api *apiInfo) Log() {
	Log(api.method, api.path, api.title)
}

type ApiInfo struct {
	Path, Method, Title string
	Version             int
	Changelog           []changelog
	Createlog           changelog
	Deprecated          *changelog
	Middleware          []http.HandlerFunc
}

func (api *ApiInfo) GetApiInfo() *apiInfo {
	return (*apiInfo)(unsafe.Pointer(api))
}
func (api *apiInfo) GetApiInfo() *ApiInfo {
	return (*ApiInfo)(unsafe.Pointer(api))
}

// 获取负责人
func (api *apiInfo) getPrincipal() string {
	if len(api.changelog) == 0 {
		return api.createlog.auth
	}
	if api.deprecated != nil {
		return api.deprecated.auth
	}
	return api.changelog[len(api.changelog)-1].auth
}

// recover捕捉panic info
func GetMethodInfo(method *reflect.Method, preUrl string, httpContext reflect.Type) (info *apiInfo) {
	if method.Name == "Service" || method.Name == "FiberService" {
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				//_,_, info.version = ParseMethodName(method.Name)
				if v.version == 0 {
					v.version = 1
				}
				v.path = preUrl + v.path
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
	/*	var err error
		defer func() {
			if err != nil {
				log.Debugf("未注册: %s 原因:%v", method.Name, err)
			}
		}()*/

	if numIn != 3 {
		//err = errors.New("method参数必须为两个")
		return
	}
	if numOut != 2 {
		//err = errors.New("method返回值必须为两个")
		return
	}
	if !methodType.In(1).ConvertibleTo(httpContext) {
		//err = errors.New("service第一个参数必须为*httpctx.Context类型")
		return
	}
	if !methodType.Out(1).Implements(ErrorType) {
		//err = errors.New("service第二个返回值必须为error类型")
		return
	}
	params := make([]reflect.Value, numIn, numIn)
	for i := 0; i < numIn; i++ {
		params[i] = reflect.New(methodType.In(i).Elem())
	}
	methodValue.Call(params)
	return nil
}

func (api *apiInfo) Swagger(doc *openapi3.T, methodType reflect.Type, tag, dec string) {
	if doc.Components.Schemas == nil {
		doc.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	}

	var pathItem *openapi3.PathItem
	if doc.Paths != nil {
		if path := doc.Paths.Value(api.path); path != nil {
			pathItem = path
		} else {
			pathItem = new(openapi3.PathItem)
		}
	} else {
		doc.Paths = openapi3.NewPaths()
		pathItem = new(openapi3.PathItem)
	}

	//我觉得路径参数并没有那么值得非用不可
	parameters := make([]*openapi3.ParameterRef, 0)
	numIn := methodType.NumIn()

	if numIn == 3 {
		if api.method == http.MethodGet {
			InType := methodType.In(2).Elem()
			for j := 0; j < InType.NumField(); j++ {
				param := &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name: InType.Field(j).Name,
						In:   "query",
					},
				}
				parameters = append(parameters, param)
			}
		} else {
			reqName := methodType.In(2).Elem().Name()
			param := &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name: reqName,
					In:   "body",
				},
			}

			param.Value.Schema = new(openapi3.SchemaRef)
			param.Value.Schema.Ref = "#/definitions/" + reqName
			parameters = append(parameters, param)
			apidoc.DefinitionsApi(param.Value.Schema.Value, reflect.New(methodType.In(2)).Elem().Interface())
		}
	}

	if !methodType.Out(0).Implements(ErrorType) {
		responses := openapi3.NewResponses()
		response := responses.Default()
		response.Ref = "#/definitions/" + methodType.Out(0).Elem().Name()
		var schema openapi3.Schema
		response.Value.Content = openapi3.NewContentWithJSONSchema(&schema)
		apidoc.DefinitionsApi(&schema, reflect.New(methodType.Out(0)).Elem().Interface())
		op := openapi3.Operation{
			Summary:     api.title,
			OperationID: api.path + api.method,
			Parameters:  parameters,
			Responses:   responses,
		}

		var tags, desc []string
		tags = append(tags, tag, api.createlog.version)
		desc = append(desc, dec, api.createlog.log)
		for i := range api.changelog {
			tags = append(tags, api.changelog[i].version)
			desc = append(desc, api.changelog[i].log)
		}
		op.Tags = tags
		op.Description = strings.Join(desc, "\n")

		switch api.method {
		case http.MethodGet:
			pathItem.Get = &op
		case http.MethodPost:
			pathItem.Post = &op
		case http.MethodPut:
			pathItem.Put = &op
		case http.MethodDelete:
			pathItem.Delete = &op
		case http.MethodOptions:
			pathItem.Options = &op
		case http.MethodPatch:
			pathItem.Patch = &op
		case http.MethodHead:
			pathItem.Head = &op
		}
	}

	doc.Paths.Set(api.path, pathItem)
}

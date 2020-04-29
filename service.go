package pick

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/liov/pick/openapi"
	"github.com/liov/pick/utils"
)

type Service interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []http.HandlerFunc)
}

var svcs = make([]Service, 0)

func RegisterService(svc ...Service) {
	svcs = append(svcs, svc...)
}

func registered() {
	isRegistered = true
	svcs = nil
}

func NewRouter(genApi bool, modName string) *Router {
	router := &Router{
		route: make(map[string][]*methodHandle),
	}
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		if preUrl[len(preUrl)-1] != '/' {
			preUrl += "/"
		}
		router.group = appendGroupSorted(router.group, groupMiddle{preUrl, middleware})
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Type.NumIn() < 2 || method.Type.NumOut() != 2 {
				continue
			}
			methodInfo := getMethodInfo(value.Method(j))
			if methodInfo == nil {
				log.Fatalf("%s未注册", method.Name)
			}
			methodInfo.path, methodInfo.version = parseMethodName(method.Name)
			if methodInfo.deprecated != nil {
				methodInfo.title += "(废弃)"
			}
			methodInfo.path = preUrl + "v" + strconv.Itoa(methodInfo.version) + "/" + methodInfo.path
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			if mh, ok := router.route[methodInfo.path]; ok {
				if _, h2, _ := getHandle(methodInfo.method, mh); h2.IsValid() {
					panic("url：" + methodInfo.path + "已注册")
				} else {
					mh = append(mh, &methodHandle{methodInfo.method, methodInfo.middleware, nil, value.Method(j)})
					router.route[methodInfo.path] = mh
				}
			} else {
				router.route[methodInfo.path] = []*methodHandle{{methodInfo.method, methodInfo.middleware, nil, value.Method(j)}}
			}
			fmt.Printf(" %s\t %s %s\t %s\n",
				utils.Green("API:"),
				utils.Yellow(utils.FormatLen(methodInfo.method, 6)),
				utils.Blue(utils.FormatLen(methodInfo.path, 50)), utils.Purple(methodInfo.title))
			if genApi {
				methodInfo.Api(value.Method(j).Type(), describe, value.Type().Name())
				openapi.WriteToFile(openapi.FilePath, modName)
			}
		}
	}
	if genApi {
		doc := GenDoc(modName, svcs)
		router.Handle(http.MethodGet, "/api-doc/md", "api文档", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(doc))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		})
	}
	registered()
	return router
}

func getMethodInfo(fv reflect.Value) (info *apiInfo) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				info = v
			} else {
				log.Panic(err)
			}
		}
	}()
	methodType := fv.Type()
	params := make([]reflect.Value, 0, fv.Type().NumIn())
	numIn := methodType.NumIn()
	numOut := methodType.NumOut()
	if numIn == 1 {
		panic("method至少一个参数且参数必须实现Session接口")
	}
	if numIn > 2 {
		panic("method参数最多为两个")
	}
	if numOut != 2 {
		panic("method返回值必须为两个")
	}
	if !methodType.In(0).Implements(contextType) {
		panic("service第一个参数必须实现Session接口")
	}
	if !methodType.Out(1).Implements(errorType) {
		panic("service第二个返回值必须为error类型")
	}
	for i := 0; i < numIn; i++ {
		params = append(params, reflect.New(methodType.In(i).Elem()))
	}
	fv.Call(params)
	return nil
}

// 从方法名称分析出接口名和版本号
func parseMethodName(originName string) (name string, version int) {
	idx := strings.LastIndexByte(originName, 'V')
	version = 1
	if idx > 0 {
		if v, err := strconv.Atoi(originName[idx+1:]); err == nil {
			version = v
		}
	} else {
		idx = len(originName)
	}
	name = utils.ConvertToSnackCase(originName[:idx])
	return
}

package pickrouter

import (
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
	"github.com/hopeio/pick"
	"net/http"
	"reflect"
)

func register(router *Router, svc ...pick.Service[http.HandlerFunc]) {
	methods := make(map[string]struct{})
	openApi(router)
	Svcs = append(Svcs, svc...)
	for _, v := range Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, HttpContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}
			methodInfoExport := methodInfo.GetApiInfo()
			router.Handle(methodInfoExport.Method, methodInfoExport.Path, methodInfoExport.Middleware, value.Method(j))
			methods[methodInfoExport.Method] = struct{}{}
			pick.Log(methodInfoExport.Method, methodInfoExport.Path, describe+":"+methodInfoExport.Title)
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.GroupApiInfos = append(pick.GroupApiInfos, &pick.GroupApiInfo{Describe: describe, Infos: infos})
		router.GroupUse(preUrl, middleware...)
	}

	allowed := make([]string, 0, 9)
	for k := range methods {
		allowed = append(allowed, k)
	}
	router.globalAllowed = allowedMethod(allowed)

	pick.Registered(Svcs)
}

func openApi(mux *Router) {
	mux.Handler(http.MethodGet, apidoc.UriPrefix+"/markdown", apidoc.Markdown)
	mux.Handler(http.MethodGet, apidoc.UriPrefix, pick.DocList)
	mux.Handler(http.MethodGet, apidoc.UriPrefix+"/swagger/*file", apidoc.Swagger)
}

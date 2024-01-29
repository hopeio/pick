package router

import (
	"github.com/hopeio/lemon/utils/log"
	"github.com/hopeio/lemon/utils/net/http/api/apidoc"
	"github.com/hopeio/pick"
	"mime"
	"net/http"
	"reflect"
)

func register(router *Router, genApiDoc bool, modName string) {
	methods := make(map[string]struct{})
	for _, v := range pick.Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, pick.ClaimsType)
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
	if genApiDoc {
		openApi(router, apidoc.FilePath, modName)
	}
	allowed := make([]string, 0, 9)
	for k := range methods {
		allowed = append(allowed, k)
	}
	router.globalAllowed = allowedMethod(allowed)

	pick.Registered()
}

func openApi(mux *Router, filePath, modName string) {
	apidoc.FilePath = filePath
	pick.Markdown(filePath, modName)
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	mux.Handler(http.MethodGet, apidoc.PrefixUri+"Markdown/*file", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filePath+".apidoc.Markdown")
	})
	pick.Swagger(filePath, modName)
	mux.Handler(http.MethodGet, apidoc.PrefixUri[:len(apidoc.PrefixUri)-1], apidoc.ApiMod)
	mux.Handler(http.MethodGet, apidoc.PrefixUri+"Swagger/*file", apidoc.HttpHandle)
}

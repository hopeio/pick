package pick

import (
	"path/filepath"

	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
)

func Swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filepath.Join(filePath+modName, modName+apidoc.SwaggerEXT))
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			methodInfo.ApiInfo.Swagger(doc, methodInfo.Method, groupApiInfo.Describe, methodInfo.Method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

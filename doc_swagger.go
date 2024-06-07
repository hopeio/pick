package pick

import (
	"github.com/hopeio/cherry/utils/net/http/apidoc"
)

func Swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filePath, modName)
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			methodInfo.ApiInfo.Swagger(doc, methodInfo.Method, groupApiInfo.Describe, methodInfo.Method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

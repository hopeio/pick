package pick

import (
	"path/filepath"

	"github.com/hopeio/lemon/utils/net/http/api/apidoc"
)

func Swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filepath.Join(filePath+modName, modName+apidoc.EXT))
	for _, groupApiInfo := range GroupApiInfos {
		for _, methodInfo := range groupApiInfo.Infos {
			methodInfo.ApiInfo.Swagger(doc, methodInfo.Method, groupApiInfo.Describe, methodInfo.Method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

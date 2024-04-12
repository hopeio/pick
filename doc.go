package pick

import "github.com/hopeio/cherry/utils/net/http/api/apidoc"

func GenApiDoc(modName string) {
	filePath := apidoc.FilePath
	Markdown(filePath, modName)
	Swagger(filePath, modName)
}

package pick

import (
	"github.com/hopeio/cherry/utils/net/http/apidoc"
	"net/http"
)

func DocList(w http.ResponseWriter, r *http.Request) {
	modName := r.URL.Query().Get("modName")
	if modName == "" {
		modName = "api"
	}
	Markdown(apidoc.ApiDocDir, modName)
	Swagger(apidoc.ApiDocDir, modName)
	apidoc.DocList(w, r)
}

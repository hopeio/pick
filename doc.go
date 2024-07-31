package pick

import (
	"github.com/hopeio/utils/net/http/apidoc"
	"net/http"
)

func DocList(w http.ResponseWriter, r *http.Request) {
	modName := r.URL.Query().Get("modName")
	if modName == "" {
		modName = "api"
	}
	Markdown(apidoc.Dir, modName)
	Swagger(apidoc.Dir, modName)
	apidoc.DocList(w, r)
}

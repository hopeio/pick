package openapi

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"gopkg.in/yaml.v2"
)

var Doc *spec.Swagger

//参数为路径和格式
func GetDoc(args ...string) *spec.Swagger {
	if Doc != nil {
		return Doc
	}
	return generate()
}

func generate() *spec.Swagger {
	Doc = new(spec.Swagger)
	info := new(spec.Info)
	Doc.Info = info

	Doc.Swagger = "2.0"
	Doc.Paths = new(spec.Paths)
	Doc.Definitions = make(spec.Definitions)

	info.Title = "Title"
	info.Description = "Description"
	info.Version = "0.01"
	info.TermsOfService = "TermsOfService"

	var contact spec.ContactInfo
	contact.Name = "Contact Name"
	contact.Email = "Contact Mail"
	contact.URL = "Contact URL"
	info.Contact = &contact

	var license spec.License
	license.Name = "License Name"
	license.URL = "License URL"
	info.License = &license

	Doc.Host = "localhost:80"
	Doc.BasePath = "/"
	Doc.Schemes = []string{"http", "https"}
	Doc.Consumes = []string{"application/json"}
	Doc.Produces = []string{"application/json"}
	return Doc
}

func WriteToFile(args ...string) {
	if Doc == nil {
		generate()
	}
	realPath := "."
	if len(args) > 0 {
		realPath = args[0]
	}

	mod := ""
	if len(args) > 1 {
		mod = args[1]
		realPath = realPath + mod
	}

	err := os.MkdirAll(realPath, 0666)
	if err != nil {
		log.Println(err)
	}

	apiType := "json"
	if len(args) > 2 {
		apiType = args[1]
	}

	realPath = filepath.Join(realPath, mod+"swagger."+apiType)

	if _, err := os.Stat(realPath); err == nil {
		os.Remove(realPath)
	}
	var file *os.File
	file, err = os.Create(realPath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	if apiType == "json" {
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(Doc)
		if err != nil {
			log.Println(err)
		}
	} else {
		b, err := yaml.Marshal(swag.ToDynamicJSON(Doc))
		if err != nil {
			log.Println(err)
		}
		if _, err := file.Write(b); err != nil {
			log.Println(err)
		}
	}
	Doc = nil
}

func NilDoc() {
	Doc = nil
}

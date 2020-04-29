package middleware

import (
	"log"
	"net/http"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Println("请求", r.RequestURI)
}

package main

import (
	"github.com/hopeio/pick/_example/router/service"
	router2 "github.com/hopeio/pick/router"
	"log"
	"net/http"

	_ "github.com/hopeio/pick/_example/router/service"
)

func main() {
	router := router2.New(&service.UserService{}, &service.TestService{}, &service.StaticService{})
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

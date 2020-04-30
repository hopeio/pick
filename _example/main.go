package main

import (
	"log"
	"net/http"

	"github.com/liov/pick"
	"github.com/liov/pick/_example/service"
)

func init() {
	pick.RegisterService(&service.UserService{})
}

func main() {
	router := pick.NewRouter(true)
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

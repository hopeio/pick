/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"log"
	"net/http"

	"github.com/hopeio/pick"
	pickstd "github.com/hopeio/pick/std"
)

func main() {
	pickstd.Register(http.DefaultServeMux, &UserService{})
	pick.OpenApi(":8081")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}

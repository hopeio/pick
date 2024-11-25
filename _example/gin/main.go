/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/gin-gonic/gin"
	"log"

	pickgin "github.com/hopeio/pick/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	pickgin.Register(server, &UserService{})
	log.Println("visit http://localhost:8080")
	log.Fatal(server.Run(":8080"))
}

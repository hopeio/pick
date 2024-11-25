/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Log(ctx *gin.Context) {
	log.Println(ctx.Request.Method, ctx.Request.RequestURI)
}

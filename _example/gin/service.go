/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/pick"
)

type UserService struct{}

func (*UserService) Service() (string, string, []gin.HandlerFunc) {
	return "用户相关", "/api/v1/user", []gin.HandlerFunc{Log}
}

type Object struct {
	Id int `json:"id,omitempty"`
}

type User struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Gender int    `json:"gender,omitempty"`
}

func (*UserService) Get(ctx *ginctx.Context, req *Object) (*User, error) {
	pick.Api(func() {
		pick.Get("/:id").
			Title("用户详情").
			CreateLog("1.0.0", "jyb", "2024/04/16", "创建").End()
	})
	// dao
	return &User{
		Id:     req.Id,
		Name:   "test",
		Gender: 1,
	}, nil
}

type Req struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

func (*UserService) GetErr(ctx *ginctx.Context, req *Req) (*User, *pick.ErrRep) {
	pick.Api(func() {
		pick.Get("/err/:id").
			Title("用户详情返回错误").
			CreateLog("1.0.0", "jyb", "2024/04/16", "创建").End()
	})
	fmt.Println(req.Name)
	// dao
	return nil, &pick.ErrRep{
		Code: 1,
		Msg:  "error",
	}
}

func (*UserService) GrpcGateway(ctx context.Context, req *Req) (*User, error) {
	pick.Api(func() {
		pick.Get("/grpcGateway").
			Title("grpcGateway").
			CreateLog("1.0.0", "jyb", "2024/04/16", "创建").End()
	})
	fmt.Println(req.Name)
	// dao
	return &User{
		Id:     req.Id,
		Name:   "test",
		Gender: 1,
	}, nil
}

func (*UserService) Middleware(ctx context.Context, req *Req) (*User, *pick.ErrRep) {
	pick.Api(func() {
		pick.Middleware(func(ctx *gin.Context) {
			fmt.Println("middleware")
		}).Get("/middleware").
			Title("中间件").
			CreateLog("1.0.0", "jyb", "2024/04/16", "创建").End()
	})
	fmt.Println(req.Name)
	// dao
	return &User{
		Id:     req.Id,
		Name:   "test",
		Gender: 1,
	}, nil
}

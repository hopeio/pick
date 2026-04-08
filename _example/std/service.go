/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hopeio/pick"
	pickstd "github.com/hopeio/pick/std"
)

type UserService struct{}

func (*UserService) Service() (string, string, []pickstd.Middleware) {
	return "用户相关", "/api/user", []pickstd.Middleware{Log, Log2, Log3}
}

type Object struct {
	Id int `uri:"id" json:"id,omitempty"`
}

type User struct {
	Id     int    `uri:"id" json:"id,omitempty" query:"id" header:"id" form:"id"`
	Name   string `json:"name,omitempty" query:"name" header:"name" form:"name"	`
	Gender int    `json:"gender,omitempty" query:"gender" header:"gender" form:"gender"`
}

func (*UserService) Uri(ctx *pickstd.Context, req *Object) (*User, error) {
	pick.Api(func() { pick.Get("/{id}").Title("用户详情").End() })
	log.Println(req.Id)
	// dao
	return &User{
		Id:     req.Id,
		Name:   time.Now().String(),
		Gender: 1,
	}, nil
}

func (*UserService) Query(ctx *pickstd.Context, req *User) (*User, error) {
	pick.Api(func() { pick.Get("").Title("用户详情").End() })
	log.Println(req.Id)
	// dao
	return req, nil
}

func (*UserService) Body(ctx *pickstd.Context, req *User) (*User, error) {
	pick.Api(func() { pick.Post("").Title("用户详情").End() })
	log.Println(req.Id)
	// dao
	return req, nil
}


func (*UserService) Header(ctx *pickstd.Context, req *User) (*User, error) {
	pick.Api(func() { pick.Put("").Title("用户详情").End() })
	log.Println(req.Id)
	// dao
	return req, nil
}


type Req struct {
	Id   int    `uri:"id" json:"id,omitempty"`
	Name string `json:"name"`
}

func (*UserService) GetErr(ctx *pickstd.Context, req *Req) (*User, *pick.ErrResp) {
	pick.Api(func() { pick.Get("/err/{id}").Title("用户详情返回错误").End() })
	fmt.Println(req.Name)
	// dao
	return nil, &pick.ErrResp{
		Code: 1,
		Msg:  "error",
	}
}

func (*UserService) GrpcGateway(ctx context.Context, req *Req) (*User, error) {
	pick.Api(func() { pick.Get("/grpcGateway").Title("grpcGateway").End() })
	fmt.Println(req.Name)
	// dao
	return &User{
		Id:     req.Id,
		Name:   time.Now().String(),
		Gender: 1,
	}, nil
}

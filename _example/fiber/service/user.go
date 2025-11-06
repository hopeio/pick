/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package service

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/fiberctx"
	"github.com/hopeio/pick"
	"github.com/hopeio/pick/_example/fiber/middle"
)

type UserService struct{}

func (*UserService) Service() (string, string, []fiber.Handler) {
	return "用户相关", "/api/v1/user", []fiber.Handler{middle.Log}
}

func (*UserService) Add(ctx *fiberctx.Context, req *SignupReq) (*TinyRep, *pick.ErrRep) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() { pick.Post("").Title("用户注册").End() })

	return &TinyRep{Msg: "测试"}, nil
}

type EditReq struct {
}
type EditReq_EditDetail struct {
}

func (*UserService) Edit(ctx *fiberctx.Context, req *EditReq) (*EditReq_EditDetail, *pick.ErrRep) {
	pick.Api(func() { pick.Put("/:id").Title("用户编辑").End() })
	return nil, nil
}

type Object struct {
	Id uint64 `json:"id"`
}

func (*UserService) Get(ctx *fiberctx.Context, req *Object) (*TinyRep, *pick.ErrRep) {
	pick.Api(func() { pick.Get("/:id").Title("用户详情").End() })

	return &TinyRep{Code: uint32(req.Id), Msg: "测试"}, nil
}

type TinyRep struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

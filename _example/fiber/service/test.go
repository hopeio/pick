/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package service

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/fiberctx"
	pick2 "github.com/hopeio/pick"
)

type TestService struct{}

func (*TestService) Service() (string, string, []fiber.Handler) {
	return "测试相关", "/api/v2/test", nil
}

type SignupReq struct {
	Mail string `json:"mail"`
}

func (*TestService) Test(ctx *fiberctx.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("").
			Title("测试").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Msg: "测试"}, nil
}

func (*TestService) Test1(ctx *fiberctx.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/").
			Title("测试1").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			End()
	})

	return &TinyRep{Msg: "测试"}, nil
}

func (*TestService) Test2(ctx *fiberctx.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/a/").
			Title("测试2").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Msg: "测试"}, nil
}

func (*TestService) Test3(ctx *fiberctx.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/a/:b").
			Title("测试3").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Msg: "测试"}, nil
}

package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/pick"
	errorsi "github.com/hopeio/utils/errors/errcode"
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

func (*UserService) GetErr(ctx *ginctx.Context, req *Req) (*User, error) {
	pick.Api(func() {
		pick.Post("/err/:id").
			Title("用户详情返回错误").
			CreateLog("1.0.0", "jyb", "2024/04/16", "创建").End()
	})
	fmt.Println(req.Name)
	// dao
	return nil, &errorsi.ErrRep{
		Code: 1,
		Msg:  "error",
	}
}

func (*UserService) GrpcGateway(ctx context.Context, req *Req) (*User, error) {
	pick.Api(func() {
		pick.Get("/grpcGateway").
			Title("用户详情返回错误").
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

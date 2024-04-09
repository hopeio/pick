package service

import (
	pick2 "github.com/hopeio/pick"
	"github.com/hopeio/tiga/context/http_context"
	"net/http"
)

type TestService struct{}

func (*TestService) Service() (string, string, []http.HandlerFunc) {
	return "测试相关", "/api/v2/test", nil
}

type SignupReq struct {
	Mail string `json:"mail"`
}

func (*TestService) Test(ctx *http_context.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("").
			Title("测试").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Message: "测试"}, nil
}

func (*TestService) Test1(ctx *http_context.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/").
			Title("测试1").
			Version(1).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			End()
	})

	return &TinyRep{Message: "测试"}, nil
}

func (*TestService) Test2(ctx *http_context.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/a/").
			Title("测试2").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Message: "测试"}, nil
}

func (*TestService) Test3(ctx *http_context.Context, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Post("/a/:b").
			Title("测试3").
			Version(3).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Message: "测试"}, nil
}

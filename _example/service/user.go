package service

import (
	"net/http"

	"github.com/liov/pick"
	"github.com/liov/pick/_example/middleware"
	"github.com/liov/pick/_example/model"
)

type UserService struct{}

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middleware.Log}
}

func (*UserService) Add(ctx *model.Claims, req *model.SignupReq) (*model.User, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Method(http.MethodPost).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})

	return &model.User{Name: "测试"}, nil
}

func (*UserService) Edit(ctx *model.Claims, req *model.User) (*model.User, error) {
	pick.Api(func() interface{} {
		return pick.Method(http.MethodPut).
			Title("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除")
	})

	return nil, nil
}

func (*UserService) LoginV2(ctx *model.Claims, req *model.LoginReq) (*model.LoginRep, error) {
	pick.Api(func() interface{} {
		return pick.Method(http.MethodGet).
			Title("用户登录").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.LoginRep{Token: ctx.Id}, nil
}

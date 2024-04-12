## [pick框架](https://github.com/actliboy/pick)
一个基于反射的自动注入api开发框架http api服务器,灵感来自于grpc和springmvc,pick的路由基于httprouter改造,如不想使用，pick同时兼容gin,fiber(fasthttp),底层可选择。


## 特点

- 摆脱路由注册
    >摆脱零散的各处的路由注册 `xxx.Handle("/"，func(){})`
- 摆脱w,r,摆脱xxx.Context
    >不再编写这样`func(w http.ResponseWriter, r *http.Request)`或者`func(ctx xxx.Context){ctx.XXX()}`的业务代码
- 专注于业务
```go
func (*UserService) Add(ctx *http_context.Context, req *model.SignupReq) (*model.User, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() {
		return pick.Post("").//定义请求的方法及路由
			Title("用户注册").//接口描述
            Middleware(nil).//中间件
            //接口迭代信息
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").//创建，唯一
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").//变更，可有多个
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除").//废弃，唯一
            End()
	})

	return &model.User{Name: "测试"}, nil
}

```  
## 使用

go get github.com/hopeio/pick

## 快速开始

首先我们需要定义一个服务
```go
type UserService struct{}
//需要实现Service方法，返回该服务的说明，url前缀，以及需要的中间件
func (*UserService) Service() (string, string, []http.HandlerFunc) {
    return "用户相关", "/api/v1/user", []http.HandlerFunc{middleware.Log}
}

```
然后可以写我们的业务方法
```go

func (*UserService) AddV2(ctx *http_context.Context, req *model.SignupReq) (*model.User, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() {
		return pick.Post("").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").
            End()
	})

	return &model.User{Name: "测试"}, nil
}


func (*UserService) Edit(ctx *http_context.Context, req *model.User) (*model.User, error) {
	pick.Api(func() {
		return pick.Put("/:id").
			Title("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除").
            End()
	})

	return nil, nil
}

```

这会生成如下的Api
```shell
 API:	 POST   /api/v1 /user   用户注册
 API:	 PUT    /api/v1/user/:id   用户编辑(废弃)
```

第一个参数可以用于身份验证，当然你必须自己实现
作为一个例子,jwt验证，这会直接在传入我们的业务方法前调用
```go
import (
  "errors"
  "github.com/golang-jwt/jwt/v5"
  "github.com/hopeio/cherry/context/http_context"
)

type AuthInfo struct {
  Id int
  jwt.RegisteredClaims
}

func ParseAuthInfo(ctx *http_context.Context) (*AuthInfo, error) {
  token := ctx.Request.Header.Get("Authorization")
  authInfo := &AuthInfo{}
  if token == "" {
    return nil, errors.New("未登录")
  }
  tokenClaims, _ := (&jwt.Parser{}).ParseWithClaims(token, authInfo, func(token *jwt.Token) (interface{}, error) {
    return "TokenSecret", nil
  })
  
  if tokenClaims != nil {
    if claims, ok := tokenClaims.Claims.(*AuthInfo); ok && tokenClaims.Valid {
    return claims,nil
    }
  }
  return nil,errors.New("未登录")
}

```
第二个参数为我们定义的参数，
返回值分别是自定义的返回结构体及错误，

> 作为一个建议这部分可以用protobuf定义，简直完美，当然这样得为第一个参数实现Context

如果是复杂场景呢,我们也是可以注册http.HandlerFunc的
```go
router.Handle(http.MethodGet, "/api-doc/md", "api文档", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(doc))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		})
```
然后，注册我们的服务

```go
func init(){
    pickrouter.RegisterService(&service.UserService{})
}
```
当然你可以注册多个
```go
func init(){
    pickrouter.RegisterService(&service.UserService{},&other.Service{})
}
```
最后启动我们的服务
```go
func main() {
    //是否生成文档
	router := pickrouter.NewRouter(true)
    pickrouter.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
```
它是如此的简单

## 文档生成

是的，你看到了`pick.NewRouter(true)`,当传入true时，会为我们生成文档
当然，这需要你定义的请求配合，例如
```go
type User struct {
	Id uint64 `json:"id"`
	Name string `json:"name" comment:"名字" validate:"gte=3,lte=10"`
	Password string `json:"password" comment:"密码" validate:"gte=8,lte=15"`
	Mail string `json:"mail" comment:"邮箱" validate:"email"`
	Phone string `json:"phone" comment:"手机" validate:"phone"`
}
```
如果你需要markdown文档，/api-doc/markdown
它会为我们生成如下文档
> [TOC]
> 
> # 用户相关  
> ----------
> ## 用户注册-v1(`/api/v1/user`)  
> **POST** `/api/v1/user` _(Principal jyb)_  
> ### 接口记录  
> |版本|操作|时间|负责人|日志|  
> | :----: | :----: | :----: | :----: | :----: |  
> |1.0.0|创建|2019/12/16|jyb|创建|  
> |1.0.1|变更|2019/12/16|jyb|修改测试|  
> ### 参数信息  
> |字段名称|字段类型|字段描述|校验要求|  
> | :----  | :----: | :----: | :----: |  
> |name|string|名字|长度必须至少为3个字符|  
> |password|string|密码|长度必须至少为8个字符|  
> |mail|string|邮箱|必须是一个有效的手机号!|  
> __请求示例__  
> ```json  
> {
> 	"name": "耰塧囎飿段",
> 	"password": "虱鷅磷黽楑",
> 	"mail": "盬艦潦昊譙"
> }  
> ```  
> ### 返回信息  
> |字段名称|字段类型|字段描述|  
> | :----  | :----: | :----: | 
> |id|number||  
> |name|string|名字|  
> |password|string|密码|  
> |mail|string|邮箱|  
> |mail|string|手机|  
> __返回示例__  
> ```json  
> {
> 	"id": 1357,
> 	"name": "鐷嚅凮珘緻",
> 	"password": "梊朖迍髽栳"
> }  
> ```  
> ## ~~用户编辑-v1(废弃)(`/api/v1/user/:id`)~~  
> **PUT** `/api/v1/user/:id` _(Principal jyb)_  
> ### 接口记录  
> ...

是的，示例并不那么好看，并非不能支持简体字和英文字母，我计划单独写一个mock模块

当然你钟情swagger，也可以`/api-doc/swagger`
![Image text](./static/1712546925271.jpg)


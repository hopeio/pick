# [pick](https://github.com/actliboy/pick)
一个基于反射的自动注入api开发框架,灵感来自于grpc和springmvc
pick的底层是灵活的,默认基于gin,同时兼容fiber(fasthttp)。

# feature

- 摆脱路由注册
    >❌`xxx.Handle("/"，func(){})`
- 摆脱w,r,摆脱xxx.Context这种不直观输入输出的handler
    >❌`func(w http.ResponseWriter, r *http.Request)`或者`func(ctx xxx.Context){ctx.XXX()}`的业务代码
- 类grpc的函数签名,专注于业务
   > ✅`func(ctx *ginctx.Context,r ReqStruct) (RespStruct,*pick.ErrRep)`

# quick start
go get github.com/hopeio/pick

`go run $(go list -m -f {{.Dir}}  github.com/hopeio/pick)/_example/gin/main.go`

# usage
## main.go
```go
func main() {
  server := gin.New()
  pickgin.Register(server, &UserService{})
  log.Fatal(server.Run(":8080"))
}
```
## service.go
```go
import (
    "github.com/hopeio/context/ginctx"
    "github.com/hopeio/pick"
    "github.com/gin-gonic/gin"
)
// 首先我们需要定义一个服务
type UserService struct{}
//需要实现Service方法，返回该服务的说明，url前缀，以及需要的中间件
func (*UserService) Service() (string, string, []gin.HandlerFunc) {
return "用户相关", "/api/v1/user", []gin.HandlerFunc{}
}
type Req struct{
  ID int `json:"id"`
}
type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
}
// 然后可以写我们的业务方法
func (*UserService) Get(ctx *ginctx.Context, req *Req) (*User, *pick.Er
rRep) {
//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
  pick.Api(func() {
    pick.Get(":/id").
    Title("用户详情").
    CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
    ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").
    End()
  })
  return &model.User{ID:req.ID,Name: "测试"}, nil
}

```  
这会生成如下的Api

 >API:	 Get   /api/v1/user/:id   用户详情

 >curl http://localhost:8080/api/v1/user/1  
 > 返回: `{"id":1,"name":"测试"}`  

# 文档生成
/apidoc
pick会为我们生成openapi和markdown文档
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
## openapi
![Image text](_assets/1712546925271.jpg)
## markdown
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
> | 字段名称     |字段类型| 字段描述 |     校验要求     |  
> |:---------| :----: |:----:|:------------:|  
> | name     |string|  名字  | 长度必须至少为3个字符  |  
> | password |string|  密码  | 长度必须至少为8个字符  |  
> | mail     |string|  邮箱  | 必须是一个有效的邮箱!  |  
> | phone    |string|  手机  | 必须是一个有效的手机号! |  
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

## 兼容grpc
```go
func (*UserService) Get(ctx context.Context, req *Req) (*User, error) {
//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
  pick.Api(func() {
    pick.Get(":/id").
    Title("用户详情").
    CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
    ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").
    End()
  })
  return &model.User{ID:req.ID,Name: "测试"}, nil
}
```

# changelog
1. 移除httprouter,默认gin路由

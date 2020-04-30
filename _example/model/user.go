package model

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name" annotation:"名字" validate:"gte=3,lte=10"`
	Password string `json:"password" annotation:"密码" validate:"gte=8,lte=15"`
	Mail     string `json:"mail" annotation:"邮箱" validate:"email"`
	Phone    string `json:"phone" annotation:"手机" validate:"phone"`
}

type SignupReq struct {
	Name     string `json:"name" annotation:"名字" validate:"gte=3,lte=10"`
	Password string `json:"password" annotation:"密码" validate:"gte=8,lte=15"`
	Mail     string `json:"mail" annotation:"邮箱" validate:"phone"`
}

type LoginReq struct {
	Password string `json:"password" annotation:"密码" validate:"gte=8,lte=15"`
	Mail     string `json:"mail" annotation:"邮箱" validate:"phone"`
}

type LoginRep struct {
	Token string `json:"token"`
}

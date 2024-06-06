package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/cherry/context/httpctx"
)

type AuthInfo struct {
	Id int
	jwt.RegisteredClaims
}

func ParseAuthInfo(ctx *httpctx.Context) (*AuthInfo, error) {
	token := ctx.RequestCtx.Request.Header.Get("Authorization")
	authInfo := &AuthInfo{}
	if token == "" {
		return nil, errors.New("未登录")
	}
	tokenClaims, _ := (&jwt.Parser{}).ParseWithClaims(token, authInfo, func(token *jwt.Token) (interface{}, error) {
		return "TokenSecret", nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*AuthInfo); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, errors.New("未登录")
}
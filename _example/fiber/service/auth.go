package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/context/fiberctx"
)

type AuthInfo struct {
	Id int
	jwt.RegisteredClaims
}

func ParseAuthInfo(ctx *fiberctx.Context) (*AuthInfo, error) {
	tokens := ctx.RequestCtx.GetReqHeaders()["Authorization"]
	if len(tokens) == 0 {
		return nil, errors.New("未登录")
	}
	authInfo := &AuthInfo{}
	token := tokens[0]
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

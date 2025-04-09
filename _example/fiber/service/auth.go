/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

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
	token := ctx.ReqCtx.RequestHeader().Get("Authorization")
	if token == "" {
		return nil, errors.New("未登录")
	}
	authInfo := &AuthInfo{}
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

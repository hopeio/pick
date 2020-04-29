package model

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId uint64
	jwt.StandardClaims
}

func (claims *Claims) GenerateToken() (string, error) {
	now := time.Now().Unix()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: now + int64(24*time.Hour),
		IssuedAt:  now,
		Issuer:    "hoper",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString("TokenSecret")

	return token, err
}

func (claims *Claims) Parse(req *http.Request) error {
	claims.Id = req.RequestURI
	return nil
}

/*
func (claims *Claims) Parse(req *http.Request) error {

	token := req.Header.Get("Authorization")

	if token == "" {
		return errors.New("未登录")
	}
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return "TokenSecret", nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return errors.New("登录超时")
			}
			return nil
		}
	}
	return errors.New("未登录")
}
*/

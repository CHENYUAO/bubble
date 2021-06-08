package handle

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	username string
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24

var Mysecret = []byte("fkuyjs")

func GenerateToken(username string) (string, error) {
	//设置token信息
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "chenyuao",
		},
	}
	//生成*jwt.Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//生成token字符串并返回
	return token.SignedString(Mysecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	//利用token字符串、自定义结构体和加密串生成*jwt.Token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return Mysecret, nil
	})
	if err != nil {
		return nil, err
	}
	//生成自定义结构体，类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

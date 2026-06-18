package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 密钥
var SecretKey = []byte("your-secret-key")

type CustomClaims struct {
	UserID               uint64 `json:"user_id"` //要放到token里的信息
	Role                 string `json:"role"`
	jwt.RegisteredClaims        //jwt标准字段
}

// 生成token
func GenerateToken(userID uint64, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                         //生效时间
			Issuer:    "blog",                                                 //发行者
		},
	}
	//token格式 ==> header.payload.signature
	//SigningMethodHS256决定header里的alg
	//claims决定payload
	//SignedString决定signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey) //使用密钥加密
}

// 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	//从token里解析出信息，并验证签名
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}
	//验证签名是否还有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	//无效则返回err
	return nil, jwt.ErrTokenInvalidClaims
}

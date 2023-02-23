package ginx

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

var myKey = []byte("douyin")

// GenerateToken 生成token
func GenerateToken(userId int64) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	UserClaim := UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("错误的token,err: %v", err)
	}
	return userClaim, nil
}

package pkgs

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	SIGNED_KEY = "53a01e6bd34caef997ddd24f5ee9d3e1"
)

type JwtToken struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   int64  `json:"expired_at"`
}

func CreateJWTToken(accountID int64) *JwtToken {
	expiredAt := time.Now().UnixMilli()
	claims := &jwt.StandardClaims{
		ExpiresAt: expiredAt,
		Subject:   fmt.Sprint(accountID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(SIGNED_KEY))
	return &JwtToken{
		AccessToken: ss,
		ExpiredAt:   expiredAt,
	}
}

// 解密
func ParseJWTToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNED_KEY), nil
	})
	if err != nil {
		return 0, Unauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, TokenParseError
	}
	defer func() {
		if e := recover(); err != nil {
			log.Printf("pase sub error %v", e)
		}
	}()

	sub, ok := claims["sub"].(string)
	var accountId int
	if ok {
		accountId, err = strconv.Atoi(sub)
		if err != nil {
			return 0, TokenParseError
		}
		return int64(accountId), nil
	}

	subInt, ok := claims["sub"].(float64)
	if ok {
		return int64(subInt), nil
	}

	log.Printf("sub is not int64 %v", claims["sub"])
	return 0, nil
}

func GetAuthorization(auth string) string {
	authArr := strings.Split(auth, " ")
	if len(authArr) == 2 {
		return authArr[1]
	}

	return ""
}

func GetAccountIdFromHeader(auth string) (int64, error) {
	token := GetAuthorization(auth)
	return ParseJWTToken(token)
}

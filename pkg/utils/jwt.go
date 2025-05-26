package utils

import (
	"errors"
	"gin-web-project/internal/config"
	"gin-web-project/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	jwtSecret   = []byte(config.Cfg.JWT.Secret)
	TokenPrefix = "user_token:"
)

func GenerateToken(userID uint, username string) (string, error) {
	// 计算令牌的过期时间
	expireDuration := time.Duration(config.Cfg.JWT.ExpireTime) * time.Hour
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(jwtSecret)
	//保存Token到redis
	database.RedisClient.Set(database.Ctx, TokenPrefix+strconv.Itoa(int(userID)), token, expireDuration)
	return signedString, err
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

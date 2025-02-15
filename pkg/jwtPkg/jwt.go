package jwtPkg

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/k1v4/avito_shop/internal/entity"
	"time"
)

// TODO to config
const secret = "secret"

func NewToken(user entity.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package jwtPkg

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/labstack/echo/v4"
	"strings"
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

func ExtractToken(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func ValidateTokenAndGetUserId(tokenString string) (int, error) {
	// парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	// проверяем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// извлекаем userId
		userId, okUser := claims["id"].(float64)
		if !okUser {
			return 0, fmt.Errorf("userId not found in token")
		}

		// проверяем срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, fmt.Errorf("token expired")
			}
		}

		return int(userId), nil
	}

	return 0, errors.New("invalid token")
}

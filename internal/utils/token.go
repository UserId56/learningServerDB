package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"learningServerDB/internal/models"
	"strings"
)

func JWTConfirm(secret string, tokenString string) (*models.User, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не верный метод кодирования")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, _ := claims["username"].(string)
		idFloat, _ := claims["userId"].(float64)
		user := &models.User{
			Id:       int64(idFloat),
			Username: username,
		}
		return user, true
	}
	return nil, false
}

func SplitBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("неверный формат токена авторизации")
	}
	return parts[1], nil
}

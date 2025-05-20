package repository

import (
	"errors"
	_ "time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthService struct {
	SecretKey string
}

func (j *JWTAuthService) Authenticate(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims == nil {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user id")
	}
	return int(userID), nil
}

func (j *JWTAuthService) HasPermission(userID int, resource string) bool {
	return true
}

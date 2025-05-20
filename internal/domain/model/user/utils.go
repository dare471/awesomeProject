package user

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет, соответствует ли хеш паролю
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
} 
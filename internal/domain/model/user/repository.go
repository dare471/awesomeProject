package user

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (User, error)
	FindByID(id uint) (User, error)
	FindAll() ([]User, error)
	Create(user *User) error
	UpdateToken(id uint, token string) error
	UpdateRefreshToken(id uint, refreshToken string) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) FindByEmail(email string) (User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, result.Error
	}
	return user, nil
}

func (r *RepositoryImpl) FindByID(id uint) (User, error) {
	var user User
	log.Printf("Attempting to find user with ID: %d", id)

	// Проверяем соединение
	if r.db == nil {
		return user, fmt.Errorf("database connection is nil")
	}

	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User not found with ID: %d", id)
			return user, errors.New("user not found")
		}
		log.Printf("Database error while finding user: %v", result.Error)
		return user, result.Error
	}

	log.Printf("Successfully found user: %+v", user)
	return user, nil
}

func (r *RepositoryImpl) FindAll() ([]User, error) {
	var users []User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *RepositoryImpl) UpdateToken(id uint, token string) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Update("token", token)
	return result.Error
}

func (r *RepositoryImpl) UpdateRefreshToken(id uint, refreshToken string) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Update("refresh_token", refreshToken)
	return result.Error
}

func (r *RepositoryImpl) Create(user *User) error {
	// Проверяем, существует ли пользователь с таким email
	var existingUser User
	result := r.db.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		return errors.New("user with this email already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// Создаем нового пользователя
	result = r.db.Create(user)
	return result.Error
}

package service

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = "secret"

type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         user.User `json:"user"`
}

type UserService struct {
	userRepo user.Repository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: user.NewRepository(database.GetDB()),
	}
}

func (s *UserService) Login(email, password string) (AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	if !user.CheckPasswordHash(password) {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return AuthResponse{}, err
	}

	// Обновляем токен в базе данных
	if err := s.userRepo.UpdateToken(user.ID, tokenString); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token:        tokenString,
		RefreshToken: "",
		User:         user,
	}, nil
}

func (s *UserService) GetUserByID(id uint) (user.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAll() ([]user.User, error) {
	return s.userRepo.FindAll()
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Age      int    `json:"age" binding:"required,min=18"`
	City     string `json:"city" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *UserService) CreateUser(req CreateUserRequest) (user.User, error) {
	// Проверяем, существует ли пользователь с таким email
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return user.User{}, errors.New("user with this email already exists")
	}

	// Создаем нового пользователя
	newUser := &user.User{
		Email:         req.Email,
		Password:      req.Password,
		Name:          req.Name,
		Age:           req.Age,
		City:          req.City,
		IsActive:      true,
		IsActive_at:   time.Now(),
		IsVerified:    true,
		IsVerified_at: time.Now(),
		IsDeleted:     false,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return user.User{}, err
	}

	return *newUser, nil
}

// func GetName(id int) (model.UserStruct, error) {
// 	user, exists := model.UserProfile[id]
// 	if !exists {
// 		return user, errors.New("user not found")
// 	}
// 	return user, nil
// }

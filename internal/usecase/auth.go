package usecase

import (
	"awesomeProject/internal/domain/model/user"
)

type AuthService interface {
	Authenticate(token string) (*user.User, error)
	GetUserProfile(token string) (*user.User, error)
}

type UserUseCase struct {
	userRepo user.Repository
}

func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

func (u *UserUseCase) Authenticate(token string) (*user.User, error) {
	// Реализация аутентификации
	return nil, nil
}

func (u *UserUseCase) GetUserProfile(token string) (*user.User, error) {
	// Реализация получения профиля пользователя
	return nil, nil
}

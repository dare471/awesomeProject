package service

import (
	"awesomeProject/internal/domain/model/news"
	"awesomeProject/internal/domain/model/upload"
	"awesomeProject/internal/domain/model/user"
)

// UserServiceInterface определяет контракт для операций с пользователями
type UserServiceInterface interface {
	// CreateUser создает нового пользователя
	CreateUser(req CreateUserRequest) (user.User, error)

	// Login выполняет аутентификацию пользователя
	Login(email, password string) (AuthResponse, error)

	// GetUserByID получает пользователя по ID
	GetUserByID(id uint) (user.User, error)

	// UpdateUser обновляет данные пользователя
	UpdateUser(id uint, user *user.User) error

	// DeleteUser удаляет пользователя
	DeleteUser(id uint) error

	// GetUserByEmail получает пользователя по email
	GetUserByEmail(email string) (user.User, error)

	// UpdateUserStatus обновляет статус пользователя
	UpdateUserStatus(id uint, isActive bool) error

	// UpdateUserVerification обновляет статус верификации пользователя
	UpdateUserVerification(id uint, isVerified bool) error
}

// AuthServiceInterface определяет контракт для операций аутентификации
type AuthServiceInterface interface {
	// Authenticate проверяет токен и возвращает пользователя
	Authenticate(token string) (*user.User, error)

	// RefreshToken обновляет токен доступа
	RefreshToken(refreshToken string) (AuthResponse, error)

	// ValidateToken проверяет валидность токена
	ValidateToken(token string) bool
}

// NewsServiceInterface определяет контракт для операций с новостями
type NewsServiceInterface interface {
	// CreateNews создает новую новость
	CreateNews(news *news.News) error

	// GetNewsByID получает новость по ID
	GetNewsByID(id uint) (news.News, error)

	// UpdateNews обновляет новость
	UpdateNews(news *news.News) error

	// DeleteNews удаляет новость
	DeleteNews(id uint) error

	// GetAllNews получает все новости
	GetAllNews() ([]news.News, error)
}

// UploadServiceInterface определяет контракт для операций с загрузками
type UploadServiceInterface interface {
	// CreateUpload создает новую загрузку
	CreateUpload(upload *upload.Upload) error

	// GetUploadByID получает загрузку по ID
	GetUploadByID(id uint) (upload.Upload, error)

	// UpdateUpload обновляет загрузку
	UpdateUpload(upload *upload.Upload) error

	// DeleteUpload удаляет загрузку
	DeleteUpload(id uint) error

	// GetAllUploads получает все загрузки
	GetAllUploads() ([]upload.Upload, error)
}

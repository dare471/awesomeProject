package service

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/config"
	"awesomeProject/internal/database"
	"awesomeProject/internal/domain/model/user"
	"awesomeProject/internal/metrics"
	"context"
	"errors"
	"fmt"
	"sync"
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
	cache    *cache.Cache
	config   *config.Config
}

// UserWithDetails содержит пользователя с дополнительными данными
type UserWithDetails struct {
	user.User
	LastLoginAt    time.Time `json:"last_login_at"`
	LoginCount     int       `json:"login_count"`
	ActiveSessions int       `json:"active_sessions"`
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: user.NewRepository(database.GetDB()),
		cache:    cache.GetCache(),
		config:   config.GetConfig(),
	}
}

func (s *UserService) Login(email, password string) (AuthResponse, error) {
	// Пробуем получить пользователя из кэша
	cacheKey := fmt.Sprintf("user:email:%s", email)
	if cachedUser, ok := s.cache.Get(cacheKey); ok {
		metrics.RecordCacheHit()
		u := cachedUser.(user.User)
		if !u.CheckPasswordHash(password) {
			return AuthResponse{}, errors.New("invalid credentials")
		}
		return s.generateAuthResponse(u)
	}
	metrics.RecordCacheMiss()

	// Если пользователя нет в кэше, получаем из базы данных
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	if !u.CheckPasswordHash(password) {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	// Сохраняем пользователя в кэш
	s.cache.Set(cacheKey, u)
	return s.generateAuthResponse(u)
}

func (s *UserService) generateAuthResponse(u user.User) (AuthResponse, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.userRepo.UpdateToken(u.ID, tokenString); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token:        tokenString,
		RefreshToken: "",
		User:         u,
	}, nil
}

func (s *UserService) GetUserByID(id uint) (user.User, error) {
	// Пробуем получить пользователя из кэша
	cacheKey := fmt.Sprintf("user:id:%d", id)
	if cachedUser, ok := s.cache.Get(cacheKey); ok {
		metrics.RecordCacheHit()
		return cachedUser.(user.User), nil
	}
	metrics.RecordCacheMiss()

	// Если пользователя нет в кэше, получаем из базы данных
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return user.User{}, err
	}

	// Сохраняем пользователя в кэш
	s.cache.Set(cacheKey, u)
	return u, nil
}

func (s *UserService) GetAll() ([]user.User, error) {
	return s.userRepo.FindAll()
}

// GetAllWithDetails возвращает всех пользователей с дополнительными данными
func (s *UserService) GetAllWithDetails(ctx context.Context) ([]UserWithDetails, error) {
	// Получаем базовый список пользователей
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make([]UserWithDetails, len(users))
	errors := make(chan error, len(users))

	// Для каждого пользователя запускаем горутину для получения дополнительных данных
	for i, u := range users {
		wg.Add(1)
		go func(index int, userItem user.User) {
			defer wg.Done()

			// Создаем каналы для получения данных
			lastLoginChan := make(chan time.Time, 1)
			loginCountChan := make(chan int, 1)
			sessionsChan := make(chan int, 1)

			// Запускаем горутины для получения разных типов данных
			go func() {
				defer close(lastLoginChan)
				// Здесь можно добавить реальную логику получения времени последнего входа
				lastLoginChan <- time.Now() // Заглушка
			}()

			go func() {
				defer close(loginCountChan)
				// Здесь можно добавить реальную логику получения количества входов
				loginCountChan <- 0 // Заглушка
			}()

			go func() {
				defer close(sessionsChan)
				// Здесь можно добавить реальную логику получения активных сессий
				sessionsChan <- 0 // Заглушка
			}()

			// Используем select для получения данных с таймаутом
			select {
			case lastLogin := <-lastLoginChan:
				select {
				case loginCount := <-loginCountChan:
					select {
					case sessions := <-sessionsChan:
						results[index] = UserWithDetails{
							User:           userItem,
							LastLoginAt:    lastLogin,
							LoginCount:     loginCount,
							ActiveSessions: sessions,
						}
					case <-time.After(s.config.Timeout.UserDetails):
						errors <- fmt.Errorf("timeout getting sessions for user %d", userItem.ID)
					}
				case <-time.After(s.config.Timeout.UserDetails):
					errors <- fmt.Errorf("timeout getting login count for user %d", userItem.ID)
				}
			case <-time.After(s.config.Timeout.UserDetails):
				errors <- fmt.Errorf("timeout getting last login for user %d", userItem.ID)
			case <-ctx.Done():
				errors <- ctx.Err()
			}
		}(i, u)
	}

	// Запускаем горутину для закрытия канала ошибок после завершения всех операций
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Собираем ошибки
	var errs []error
	for err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Если есть ошибки, возвращаем их
	if len(errs) > 0 {
		return results, fmt.Errorf("errors occurred while fetching user details: %v", errs)
	}

	return results, nil
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

	// Сохраняем пользователя в кэш
	cacheKey := fmt.Sprintf("user:email:%s", req.Email)
	s.cache.Set(cacheKey, *newUser)

	return *newUser, nil
}

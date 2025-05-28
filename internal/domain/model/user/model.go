package user

import (
	"awesomeProject/internal/domain/model/common"
	"time"

	"gorm.io/gorm"
)

// User представляет модель пользователя
type User struct {
	common.Base
	Name          string    `json:"name" gorm:"size:255;not null"`
	Age           int       `json:"age" gorm:"not null"`
	City          string    `json:"city" gorm:"size:255;not null"`
	Password      string    `json:"-" gorm:"size:255;not null"`
	Email         string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Token         string    `json:"-"`
	RefreshToken  string    `json:"refresh_token,omitempty"`
	Role          string    `json:"role" gorm:"size:255;not null"`
	IsActive      bool      `json:"is_active" gorm:"null" gorm:"default:false"`
	IsActive_at   time.Time `json:"is_active_at" gorm:"null" gorm:"default:false"`
	IsVerified    bool      `json:"is_verified" gorm:"null" gorm:"default:false"`
	IsVerified_at time.Time `json:"is_verified_at" gorm:"null" gorm:"default:false"`
	IsDeleted     bool      `json:"is_deleted" gorm:"null" gorm:"default:false"`
}

// TableName указывает имя таблицы в базе данных
func (User) TableName() string {
	return "users_struct"
}

// BeforeSave - хук GORM для хеширования пароля перед сохранением
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

// CheckPasswordHash проверяет, соответствует ли хеш паролю
func (u *User) CheckPasswordHash(password string) bool {
	return CheckPasswordHash(password, u.Password)
}

package user_deleted

import (
	"awesomeProject/internal/domain/model/common"
	"gorm.io/gorm"
)

type UserDeleted struct {
	common.Base
	user.User
	ID uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null"`
	DeletedAt time.Time `json:"deleted_at" gorm:"not null"`
}

func (UserDeleted) TableName() string {
	return "users_deleted_struct"
}



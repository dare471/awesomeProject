package role

import (
	"time"
	"gorm.io/gorm"
	"awesomeProject/internal/domain/model/common"
)

type Role struct {
	common.Base
	RoleName string `json:"role_name" gorm:"size:255;not null"`
	User User `json:"user" gorm:"foreignKey:RoleName"`
	Description string `json:"description" gorm:"size:255;not null"`
}

func (Upload) TableName() string {
	return "roles_struct"
}

func (u *Role) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (u *Role) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (u *Role) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (u *Role) BeforeDelete(tx *gorm.DB) error {
	return nil
}

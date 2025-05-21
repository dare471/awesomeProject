package role

import (
	"gorm.io/gorm"
	"awesomeProject/internal/domain/model/common"
	"awesomeProject/internal/domain/model/user"
)

type Role struct {
	common.Base
	RoleName string `json:"role_name" gorm:"type:varchar(255);not null;unique"`
	User user.User `json:"user" gorm:"foreignKey:RoleName"`
	Description string `json:"description" gorm:"type:text;not null"`
}

func (Role) TableName() string {
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

package upload

import (
	"awesomeProject/internal/domain/model/common"
	"gorm.io/gorm"
)

type Upload struct {
	common.Base
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"size:255;not null"`
	Content     string `json:"content" gorm:"size:255;not null"`
	Type        string `json:"type" gorm:"size:255;not null"`
	Path        string `json:"path" gorm:"size:255;not null"`
}

func (Upload) TableName() string {
	return "uploads_struct"
}

func (u *Upload) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (u *Upload) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (u *Upload) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (u *Upload) BeforeDelete(tx *gorm.DB) error {
	return nil
}


package news

import (
	"awesomeProject/internal/domain/model/common"
	"gorm.io/gorm"
)

type News struct {
	common.Base
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"size:255;not null"`
	Content     string `json:"content" gorm:"size:255;not null"`
	Author      string `json:"author" gorm:"size:255;not null"`
	Category    string `json:"category" gorm:"size:255;not null"`
	Image       string `json:"image" gorm:"size:255;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (News) TableName() string {
	return "news_struct"
}

func (n *News) BeforeSave(tx *gorm.DB) error {
	return nil
}
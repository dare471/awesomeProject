package common

import (
	"time"
	"gorm.io/gorm"
)

// Base содержит общие поля для всех моделей
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate - хук GORM, который вызывается перед созданием записи
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

// BeforeUpdate - хук GORM, который вызывается перед обновлением записи
func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
} 
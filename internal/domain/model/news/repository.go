package news

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(news *News) error
	FindAll() ([]News, error)
	FindByID(id uint) (News, error)
	Update(news *News) error
	Delete(id uint) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewsRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(news *News) error {
	return r.db.Create(news).Error
}

func (r *RepositoryImpl) FindAll() ([]News, error) {
	var news []News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *RepositoryImpl) FindByID(id uint) (News, error) {
	var news News
	log.Printf("Attempting to find news with ID: %d", id)

	// Проверяем соединение
	if r.db == nil {
		return news, fmt.Errorf("database connection is nil")
	}

	result := r.db.First(&news, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("News not found with ID: %d", id)
			return news, errors.New("news not found")
		}
		log.Printf("Database error while finding news: %v", result.Error)
		return News{}, result.Error
	}
	return news, nil
}

func (r *RepositoryImpl) Update(news *News) error {
	return r.db.Model(&News{}).Where("id = ?", news.ID).Updates(news).Error
}

func (r *RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&News{}, id).Error
}

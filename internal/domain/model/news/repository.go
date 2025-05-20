package news

import (
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
	if err := r.db.First(&news, id).Error; err != nil {
		return News{}, err
	}
	return news, nil
}

func (r *RepositoryImpl) Update(news *News) error {
	return r.db.Model(&News{}).Where("id = ?", news.ID).Updates(news).Error
}

func (r *RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&News{}, id).Error
}

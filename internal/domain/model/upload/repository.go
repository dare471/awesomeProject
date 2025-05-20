package upload

import (
	"awesomeProject/internal/domain/model/common"
	"gorm.io/gorm"
)

type Repository interface {
	Create(upload *Upload) error
	FindAll() ([]Upload, error)
	FindByID(id uint) (Upload, error)
	Update(upload *Upload) error
	Delete(id uint) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func (r *RepositoryImpl) Create(upload *Upload) error {
	return r.db.Create(upload).Error
}

func (r *RepositoryImpl) FindAll() ([]Upload, error) {
	return r.db.Find(&Upload{}).Error
}

func (r *RepositoryImpl) FindById(id uint) (Upload, error) {
	var upload Upload
	if err := r.db.First(&upload, id).Error; err != nil {
		return Upload{}, err
	}
	return upload, nil
}

func (r *RepositoryImpl) Update(upload *Upload) error {
	return r.db.Model(&Upload{}).Where("id = ?", upload.ID).Updates(upload).Error
}

func (r *RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&Upload{}, id).Error
}

package role

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(id uint) (Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(id uint) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) FindByID(id uint) (Role, error) {
	var role Role
	if err := r.db.First(&role, id).Error; err != nil {
		return Role{}, err
	}
	return role, nil
}

func (r *RepositoryImpl) Create(role *Role) error {
	return r.db.Create(role).Error
}

func (r *RepositoryImpl) Update(role *Role) error {
	return r.db.Save(role).Error
}

func (r *RepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&Role{}, id).Error
}
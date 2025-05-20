package user_deleted

import (
	"awesomeProject/internal/domain/model/common"
	"gorm.io/gorm"
)

type Repository interface {
	Create(userDeleted *UserDeleted) error
	FindAll() ([]UserDeleted, error)
	FindByID(id uint) (UserDeleted, error)
	Update(userDeleted *UserDeleted) error
	Restore(id uint) error
}
type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(userDeleted *UserDeleted) error {
	return r.db.Create(userDeleted).Error
}

func (r *RepositoryImpl) FindAll() ([]UserDeleted, error) {
	var userDeleted []UserDeleted
	if err := r.db.Find(&userDeleted).Error; err != nil {
		return nil, err
	}
}

func (r *RepositoryImpl) FindByID(id uint) (UserDeleted, error) {
	var userDeleted UserDeleted
	if err := r.db.First(&userDeleted, id).Error; err != nil {
		return UserDeleted{}, err
	}
}

func (r *RepositoryImpl) Update(userDeleted *UserDeleted) error {
	return r.db.Model(&UserDeleted{}).Where("id = ?", userDeleted.ID).Updates(userDeleted).Error
}

func (r *RepositoryImpl) Restore(id uint) error {
	return r.db.Model(&UserDeleted{}).Where("id = ?", id).Update("is_deleted", false).Error
}

func (r *RepositoryImpl) Create(userDeleted *UserDeleted) error {
	return r.db.Create(userDeleted).Error
}

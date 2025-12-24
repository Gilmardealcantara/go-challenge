package categories

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Category, error)
	Create(category *Category) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) Create(category *Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

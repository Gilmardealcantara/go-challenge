package products

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProducts() ([]Product, error)
	GetProductsWithPagination(offset, limit int) ([]Product, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetProductsWithPagination(offset, limit int) ([]Product, int64, error) {
	var products []Product

	result := r.db.Preload("Category").Preload("Variants").Order("code").Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return products, int64(len(products)), nil
}

package products

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Product, error)
	GetWithPagination(offset, limit int) ([]Product, error)
	GetWithFilters(offset, limit int, categoryCode string, priceLessThan *float64) ([]Product, error)
	Total() (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetWithPagination(offset, limit int) ([]Product, error) {
	var products []Product

	result := r.db.Preload("Category").Preload("Variants").Order("code").Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *repository) GetWithFilters(offset, limit int, categoryCode string, priceLessThan *float64) ([]Product, error) {
	var products []Product

	// Get paginated products
	result := r.db.Preload("Category").Preload("Variants").Order("code")

	if categoryCode != "" {
		result = result.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", categoryCode)
	}

	if priceLessThan != nil {
		result = result.Where("price < ?", *priceLessThan)
	}

	if err := result.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) Total() (int64, error) {
	var count int64
	if err := r.db.Model(&Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

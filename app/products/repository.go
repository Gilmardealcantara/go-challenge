package products

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProducts() ([]Product, error)
	GetProductsWithPagination(offset, limit int) ([]Product, int64, error)
	GetProductsWithFilters(offset, limit int, categoryCode string, priceLessThan *float64) ([]Product, int64, error)
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

func (r *repository) GetProductsWithFilters(offset, limit int, categoryCode string, priceLessThan *float64) ([]Product, int64, error) {
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
		return nil, 0, err
	}

	return products, int64(len(products)), nil
}

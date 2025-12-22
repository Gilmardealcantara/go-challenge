package catalog

import (
	"gorm.io/gorm"
)

type FilterParams struct {
	Offset        int
	Limit         int
	CategoryCode  string
	PriceLessThan *float64
}

type Repository interface {
	GetAll() ([]Product, error)
	GetWithFilters(params FilterParams) ([]Product, error)
	Total() (int64, error)
	GetByCode(code string) (*Product, error)
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

func (r *repository) GetWithFilters(params FilterParams) ([]Product, error) {
	var products []Product

	// Get paginated products
	result := r.db.Preload("Category").Preload("Variants").Order("code")

	if params.CategoryCode != "" {
		result = result.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", params.CategoryCode)
	}

	if params.PriceLessThan != nil {
		result = result.Where("price < ?", *params.PriceLessThan)
	}

	if err := result.Offset(params.Offset).Limit(params.Limit).Find(&products).Error; err != nil {
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

func (r *repository) GetByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

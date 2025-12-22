package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code, a price, and belongs to a category.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID *uint           `gorm:"index"`
	Category   *categories.Category
	Variants   []variants.Variant `gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "products"
}

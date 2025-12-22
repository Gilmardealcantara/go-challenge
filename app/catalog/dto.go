package catalog

type Response struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

type ProductResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
}

type ProductDetailsResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
	Variants []VariantResponse `json:"variants"`
}

type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

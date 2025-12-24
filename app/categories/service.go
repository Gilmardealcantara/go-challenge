package categories

type Service interface {
	GetCategories() ([]CategoryResponse, error)
	CreateCategory(req CreateCategoryRequest) (*CategoryResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) GetCategories() ([]CategoryResponse, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	return responses, nil
}

func (s *service) CreateCategory(req CreateCategoryRequest) (*CategoryResponse, error) {
	category := &Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}, nil
}

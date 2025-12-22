package categories

type Service interface {
	GetCategories() ([]CategoryResponse, error)
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

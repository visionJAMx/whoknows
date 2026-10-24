package service

import "context"

type SearchService struct {
}

type SearchInput struct {
	Query    string
	Language string
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) SearchPage(ctx context.Context, input SearchInput) (string, error) {

	// TODO: call repository layer

	return "GET /", nil
}

func (s *SearchService) Search(ctx context.Context, input SearchInput) (string, error) {

	// TODO: call repository layer

	return "GET /api/search", nil
}

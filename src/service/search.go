package service

import (
	"context"

	"github.com/visionJAMx/whoknows/src/models"
)

type PageRepository interface {
	Search(ctx context.Context, query, language string) ([]models.Page, error)
}

type SearchService struct {
	pages PageRepository
}

type SearchInput struct {
	Query    string
	Language string
}

func NewSearchService(pages PageRepository) *SearchService {
	return &SearchService{pages: pages}
}

func (s *SearchService) SearchPage(ctx context.Context, input SearchInput) (string, error) {

	// TODO: call repository layer

	return "GET /", nil
}

func (s *SearchService) Search(ctx context.Context, input SearchInput) (string, error) {

	// TODO: call repository layer

	return "GET /api/search", nil
}

package service

import (
	"context"
	"fmt"
	"strings"

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

// Search returns the pages matching the query. An empty query returns no pages without searching.
func (s *SearchService) Search(ctx context.Context, input SearchInput) ([]models.Page, error) {
	input.Query = strings.TrimSpace(input.Query)
	if input.Query == "" {
		return []models.Page{}, nil
	}

	pages, err := s.pages.Search(ctx, input.Query, input.Language)
	if err != nil {
		return nil, fmt.Errorf("search pages: %w", err)
	}

	return pages, nil
}

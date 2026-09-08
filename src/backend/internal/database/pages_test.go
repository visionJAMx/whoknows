package database

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateAndSearchPages(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(databasePath)
	if err != nil {
		t.Fatalf("Open() returned an error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := Initialize(ctx, db); err != nil {
		t.Fatalf("Initialize() returned an error: %v", err)
	}

	pages := []Page{
		{
			Title:       "Go programming",
			URL:         "https://example.com/go",
			Language:    "en",
			LastUpdated: time.Now().UTC(),
			Content:     "Learn Go and Gin",
		},
		{
			Title:       "Dansk DevOps",
			URL:         "https://example.com/devops",
			Language:    "da",
			LastUpdated: time.Now().UTC(),
			Content:     "Lær om DevOps",
		},
	}

	for _, page := range pages {
		if err := CreatePage(ctx, db, page); err != nil {
			t.Fatalf("CreatePage() returned an error: %v", err)
		}
	}

	results, err := SearchPages(ctx, db, "Go", "en")
	if err != nil {
		t.Fatalf("SearchPages() returned an error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 search result, got %d", len(results))
	}

	if results[0].Title != "Go programming" {
		t.Errorf(
			"expected title %q, got %q",
			"Go programming",
			results[0].Title,
		)
	}
}

func TestSearchPagesTreatsInputAsData(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(databasePath)
	if err != nil {
		t.Fatalf("Open() returned an error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := Initialize(ctx, db); err != nil {
		t.Fatalf("Initialize() returned an error: %v", err)
	}

	results, err := SearchPages(ctx, db, `%' OR 1=1 --`, "")
	if err != nil {
		t.Fatalf("SearchPages() returned an error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 search results, got %d", len(results))
	}
}

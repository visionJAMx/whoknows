package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/visionJAMx/whoknows/src/domain"
)

func CreatePage(
	ctx context.Context,
	db *sql.DB,
	page domain.Page,
) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO pages (title, url, language, last_updated, content)
		 VALUES (?, ?, ?, ?, ?)`,
		page.Title,
		page.URL,
		page.Language,
		page.LastUpdated,
		page.Content,
	)
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	return nil
}

func SearchPages(
	ctx context.Context,
	db *sql.DB,
	query string,
	language string,
) ([]domain.Page, error) {
	searchTerm := "%" + query + "%"

	rows, err := db.QueryContext(
		ctx,
		`SELECT title, url, language, last_updated, content
		 FROM pages
		 WHERE (? = '' OR language = ?)
		   AND (title LIKE ? OR content LIKE ?)
		 ORDER BY title`,
		language,
		language,
		searchTerm,
		searchTerm,
	)
	if err != nil {
		return nil, fmt.Errorf("search pages: %w", err)
	}
	defer rows.Close()

	var pages []domain.Page

	for rows.Next() {
		var page domain.Page

		if err := rows.Scan(
			&page.Title,
			&page.URL,
			&page.Language,
			&page.LastUpdated,
			&page.Content,
		); err != nil {
			return nil, fmt.Errorf("read search result: %w", err)
		}

		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate search results: %w", err)
	}

	return pages, nil
}

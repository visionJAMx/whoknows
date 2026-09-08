package domain

import "time"

type Page struct {
	Title       string
	URL         string
	Language    string
	LastUpdated time.Time
	Content     string
}

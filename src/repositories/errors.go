package repositories

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)

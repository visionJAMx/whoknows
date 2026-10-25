package repositories

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")

	// both wrap ErrDuplicate, so errors.Is(err, ErrDuplicate) still matches
	ErrUsernameTaken = fmt.Errorf("username %w", ErrDuplicate)
	ErrEmailTaken    = fmt.Errorf("email %w", ErrDuplicate)
)

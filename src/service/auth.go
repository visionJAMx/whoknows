package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/visionJAMx/whoknows/src/domain"
	"github.com/visionJAMx/whoknows/src/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials er bevidst generisk, så login ikke afslører eksisterende brugernavne.
var ErrInvalidCredentials = errors.New("invalid username or password")

// HashPassword laver en langsom bcrypt-hash, som registreringsflowet kan gemme sikkert.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Authenticate finder brugeren og sammenligner det indsendte password med den gemte bcrypt-hash.
func Authenticate(
	ctx context.Context,
	db *sql.DB,
	username string,
	password string,
) (domain.User, error) {
	if strings.TrimSpace(username) == "" || password == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := repository.FindUserByUsername(ctx, db, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrInvalidCredentials
		}

		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}

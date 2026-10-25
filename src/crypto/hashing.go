package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// ErrMismatchedPassword is returned by Compare when the password doesn't match the hash
var ErrMismatchedPassword = bcrypt.ErrMismatchedHashAndPassword

func HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func Compare(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err
}

package domain

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
}

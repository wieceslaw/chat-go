package auth

import "time"

type (
	UserId    int64
	AuthToken string
	UserRole  string
)

const (
	DefaultRole UserRole = "user"
)

type LoginData struct {
	Username string
	Password string
}

type RegisterUser struct {
	Name     string
	Password string
}

type NewUser struct {
	Name         string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User struct {
	Id           UserId
	Name         string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

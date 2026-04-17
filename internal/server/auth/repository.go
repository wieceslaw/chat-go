package auth

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // To register the driver.
)

var (
	ErrNameAlreadyUsed = errors.New("name is already used")
	ErrUserNotFound    = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *NewUser) error
	GetUser(ctx context.Context, username string) (*User, error)
}

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		DB: db,
	}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *NewUser) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO users
		(name, password_hash, updated_at, created_at)
		VALUES
		($1, $2, $3, $4)
	`, user.Name, user.PasswordHash, user.UpdatedAt, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) GetUser(ctx context.Context, username string) (*User, error) {
	stmt, err := r.DB.Prepare("SELECT * FROM users WHERE name = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(username).Scan(&user.Id, &user.Name, &user.PasswordHash, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// --- mock ---
func MockUserRepository() UserRepository {
	return &mockUserRepository{
		make(map[string]User),
	}
}

type mockUserRepository struct {
	users map[string]User
}

func (r *mockUserRepository) Close() error {
	return nil
}

func (mr *mockUserRepository) CreateUser(ctx context.Context, user *NewUser) error {
	mr.users[user.Name] = User{
		Id:           UserId(len(mr.users)),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	return nil
}

func (mr *mockUserRepository) GetUser(ctx context.Context, username string) (*User, error) {
	user, exists := mr.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

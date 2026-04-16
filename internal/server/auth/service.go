package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegisterParams = errors.New("invalid username or password")
	ErrInvalidLoginParams    = errors.New("invalid username or password")
	ErrInvalidToken          = errors.New("invalid token")
)

type UserService interface {
	Register(ctx context.Context, user *RegisterUser) error
	Login(ctx context.Context, loginData *LoginData) (*AuthToken, error)
	ValidateToken(ctx context.Context, token AuthToken) (*User, error)
}

type JwtSecretProvider interface {
	Secret() []byte
}

func MockJwtProvider() JwtSecretProvider {
	return &mockJwtProvider{
		secret: []byte("secret"),
	}
}

type mockJwtProvider struct {
	secret []byte
}

func (p *mockJwtProvider) Secret() []byte {
	return p.secret
}

func NewUserService(
	ctx context.Context,
	repository UserRepository,
	secretProvider JwtSecretProvider,
) (UserService, error) {
	return &userServiceImpl{
		repository:     repository,
		secretProvider: secretProvider,
	}, nil
}

type userServiceImpl struct {
	repository     UserRepository
	secretProvider JwtSecretProvider
}

func (s *userServiceImpl) Register(ctx context.Context, newUser *RegisterUser) error {
	if newUser.Name == "" || newUser.Password == "" {
		return ErrInvalidRegisterParams
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to create password hash %v", err)
	}

	now := time.Now()
	user := NewUser{
		Name:         newUser.Name,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = s.repository.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, loginData *LoginData) (*AuthToken, error) {
	if loginData.Username == "" || loginData.Password == "" {
		return nil, ErrInvalidLoginParams
	}

	user, err := s.repository.GetUser(ctx, loginData.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user %v", err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginData.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to compare hash %v", err)
	}

	return s.createAuthToken(user)
}

type customClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *userServiceImpl) createAuthToken(user *User) (*AuthToken, error) {
	claims := customClaims{
		UserID:   int(user.Id),
		Username: user.Name,
		Role:     string(DefaultRole),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "chat-app",
			Subject:   fmt.Sprintf("%d", user.Id),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretProvider.Secret())
	if err != nil {
		return nil, err
	}

	authToken := AuthToken(tokenString)

	return &authToken, nil
}

func (s *userServiceImpl) ValidateToken(ctx context.Context, authToken AuthToken) (*User, error) {
	parsedToken, err := jwt.ParseWithClaims(string(authToken), &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretProvider.Secret(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*customClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	user, err := s.repository.GetUser(ctx, claims.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

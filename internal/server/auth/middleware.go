package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
)

const userKey = "user"

type AuthMiddleware struct {
	Service UserService
}

func NewAuthMiddleware(service UserService) *AuthMiddleware {
	return &AuthMiddleware{
		Service: service,
	}
}

func GetUser(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userKey).(*User)
	if ok {
		return user, nil
	}
	return nil, errors.New("No user in context")
}

func (am *AuthMiddleware) Wrap(hanlder func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := am.Service.ValidateToken(r.Context(), AuthToken(token))

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		hanlder(w, r.WithContext(ctx))
	})
}

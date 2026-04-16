package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wieceslaw/chat-go/internal/server/auth"
	"github.com/wieceslaw/chat-go/internal/server/hello"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	repository := auth.NewUserRepository("postgresql://myuser:mypassword@localhost/mydatabase?sslmode=disable")
	defer repository.Close()
	service, _ := auth.NewUserService(ctx, repository, auth.MockJwtProvider())
	middleware := auth.NewAuthMiddleware(service)
	handler := auth.NewAuthHanlder(service)
	handler.Register(mux)

	hello.Register(mux, middleware.Wrap)

	http.ListenAndServe(":8080", mux)

	fmt.Println("Server started on port: 8080")
}

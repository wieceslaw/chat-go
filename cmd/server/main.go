package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wieceslaw/chat-go/cmd/server/auth"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	repository := auth.MockUserRepository()
	service, _ := auth.NewUserService(ctx, repository, auth.MockJwtProvider())
	handler := auth.NewAuthHanlder(service)
	handler.Register(mux)

	http.ListenAndServe(":8080", mux)

	fmt.Println("Server started on port: 8080")
}

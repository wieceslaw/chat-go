package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wieceslaw/chat-go/internal/server/auth"
	"github.com/wieceslaw/chat-go/internal/server/hello"
)

func main() {
	ctx := context.Background()

	r := gin.Default()

	repository := auth.NewUserRepository("postgresql://myuser:mypassword@localhost/mydatabase?sslmode=disable")
	defer repository.Close()
	service, _ := auth.NewUserService(ctx, repository, auth.MockJwtProvider())

	authHandler := auth.NewAuthHanlder(service)
	authHandler.RegisterRoutes(r.Group(""))

	authMiddleware := auth.NewAuthMiddleware(service)

	api := r.Group("/api/v1")
	api.Use(authMiddleware.AuthRequired())
	{
		helloHandler := hello.NewHelloHandler()
		helloHandler.RegisterRoutes(api.Group("/hello"))
	}

	fmt.Println("Server started on port: 8080")
	r.Run(":8080")
}

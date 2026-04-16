package hello

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wieceslaw/chat-go/internal/server/auth"
)

type HelloHandler struct {
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

func (h *HelloHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", h.hello)
}

type helloResponseDto struct {
	Message string
}

func (h *HelloHandler) hello(c *gin.Context) {
	user, err := auth.GetUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, helloResponseDto{Message: "Hello, " + user.Name + "!"})
}

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	service UserService
}

func NewAuthHanlder(service UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/api/v1/auth/register", h.handleRegister)
	rg.POST("/api/v1/auth/login", h.handleLogin)
}

func (h *AuthHandler) handleRegister(c *gin.Context) {
	var req registerRequestDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Register(c.Request.Context(), &RegisterUser{
		Name:     req.Username,
		Password: req.Password,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ok",
	})
}

func (h *AuthHandler) handleLogin(c *gin.Context) {
	var req loginRequestDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	token, err := h.service.Login(c.Request.Context(), &LoginData{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": string(*token),
	})
}

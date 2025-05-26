package handler

import (
	"gin-web-project/internal/model"
	"gin-web-project/internal/service"
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(db *gorm.DB, apiV1 *gin.RouterGroup) *AuthHandler {
	handler := &AuthHandler{
		authService: service.NewAuthService(db),
	}
	apiV1.POST("/login", handler.Login)
	apiV1.POST("/register", handler.Register)
	return handler
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, user)
}

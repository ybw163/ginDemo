package handler

import (
	"gin-web-project/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(db *gorm.DB, apiV1 *gin.RouterGroup) *UserHandler {
	handler := &UserHandler{
		userService: service.NewUserService(db),
	}
	admin := apiV1.Group("/admin")
	admin.GET("/users", handler.Users)
	return handler
}

func (h UserHandler) Users(context *gin.Context) {
	users := h.userService.Users()
	context.JSON(200, gin.H{
		"users": users,
	})
}

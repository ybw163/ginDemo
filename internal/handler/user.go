package handler

import (
	"gin-web-project/internal/service"
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(apiV1 *gin.RouterGroup) *UserHandler {
	handler := &UserHandler{
		userService: service.NewUserService(),
	}
	admin := apiV1.Group("/admin")
	admin.GET("/users", handler.Users)
	return handler
}

func (h UserHandler) Users(context *gin.Context) {
	users := h.userService.Users()
	utils.Success(context, users)
}

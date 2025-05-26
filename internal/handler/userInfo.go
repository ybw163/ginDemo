package handler

import (
	"gin-web-project/internal/model"
	"gin-web-project/internal/service"
	"gin-web-project/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInfoHandler struct {
	userInfoService *service.UserInfoService
}

func NewUserInfoHandler(db *gorm.DB, apiV1 *gin.RouterGroup) *UserInfoHandler {
	handler := &UserInfoHandler{
		userInfoService: service.NewUserInfoService(db),
	}
	admin := apiV1.Group("/info")
	admin.GET("/", handler.info)
	admin.POST("/add", handler.add)
	return handler
}

func (h UserInfoHandler) info(context *gin.Context) {
	value, _ := context.Get("userId")
	userId := value.(uint)
	info := h.userInfoService.Info(userId)
	utils.Success(context, info)
}

func (h UserInfoHandler) add(context *gin.Context) {
	var req model.UserInfo
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(context, err.Error())
		return
	}
	value, _ := context.Get("userId")
	userId := value.(uint)
	req.UserId = userId
	h.userInfoService.Add(req)
	utils.Success(context, nil)
}

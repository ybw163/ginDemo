package service

import (
	"gin-web-project/internal/model"
	"gorm.io/gorm"
)

type UserInfoService struct {
	db *gorm.DB
}

func NewUserInfoService(db *gorm.DB) *UserInfoService {
	return &UserInfoService{db: db}
}

func (u *UserInfoService) Info(userId uint) model.UserInfo {
	var info model.UserInfo
	u.db.Where("user_id = ?", userId).First(&info)
	return info
}

func (u *UserInfoService) Add(req model.UserInfo) {
	u.db.Create(&req)
}

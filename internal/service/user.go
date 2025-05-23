package service

import (
	"gin-web-project/internal/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) Users() []model.User {
	var users []model.User
	u.db.Find(&users)
	return users
}

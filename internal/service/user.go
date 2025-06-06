package service

import (
	"gin-web-project/internal/database"
	"gin-web-project/internal/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{db: database.Db}
}

func (u *UserService) Users() []model.User {
	var users []model.User
	u.db.Find(&users)
	return users
}

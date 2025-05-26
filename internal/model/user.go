package model

import "gin-web-project/internal/database"

type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
	Password string `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	Avatar   string `json:"avatar" gorm:"size:255"`
	Status   int    `json:"status" gorm:"default:1"` // 1:正常 0:禁用
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func init() {
	database.RegisteredModels = append(database.RegisteredModels, &User{})
}

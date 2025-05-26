package model

import "gin-web-project/internal/database"

type UserInfo struct {
	BaseModel
	UserId uint `json:"user_id" gorm:"uniqueIndex;not null"`
	Age    uint `json:"age" gorm:"uniqueIndex;size:100;not null" `
	Sex    uint `json:"sex" gorm:"uniqueIndex;size:100;not null" `
	Height uint `json:"height" gorm:"uniqueIndex;size:100;not null" `
}

func init() {
	database.RegisteredModels = append(database.RegisteredModels, &UserInfo{})
}

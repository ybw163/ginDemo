package service

import (
	"errors"
	"gin-web-project/internal/database"
	"gin-web-project/internal/model"
	"gin-web-project/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{db: database.Db}
}

func (s *AuthService) Login(username, password string) (string, error) {
	var user model.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	// 检查用户是否已存在
	var count int64
	s.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, errors.New("username or email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   1,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

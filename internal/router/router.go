package router

import (
	"gin-web-project/internal/config"
	"gin-web-project/internal/handler"
	"gin-web-project/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(gin.Recovery())

	// 健康检查
	r.GET("/health", handler.HealthCheck)
	// API路由组
	api := r.Group("/api/v1")
	// 公开路由
	public := api.Group("/")
	handler.NewAuthHandler(db, public)
	// 需要认证的路由
	auth := api.Group("/")
	auth.Use(middleware.JWTAuth())
	// 管理员路由组
	// auth.Use(middleware.AdminAuth()) // 可添加管理员权限中间件
	handler.NewUserHandler(db, auth)
	handler.NewUserInfoHandler(db, auth)

	return r
}

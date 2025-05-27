package router

import (
	"gin-web-project/internal/config"
	"gin-web-project/internal/handler"
	"gin-web-project/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 设置运行模式
	gin.SetMode(config.Cfg.Server.Mode)
	r := gin.New()
	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(gin.Recovery())
	//设置handler
	setUpHandler(r)
	return r
}

func setUpHandler(r *gin.Engine) {
	// 健康检查
	handler.HealthCheck(r)
	// API路由组
	api := r.Group("/api/v1")
	// 公开路由
	public := api.Group("/")
	handler.NewAuthHandler(public)
	// 需要认证的路由
	auth := api.Group("/")
	auth.Use(middleware.JWTAuth())
	// 管理员路由组
	// auth.Use(middleware.AdminAuth()) // 可添加管理员权限中间件
	handler.NewUserHandler(auth)
	handler.NewUserInfoHandler(auth)
}

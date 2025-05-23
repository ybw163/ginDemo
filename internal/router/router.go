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

	// 初始化处理器
	authHandler := handler.NewAuthHandler(db)
	userHandler := handler.NewUserHandler(db)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由
		public := api.Group("/")
		{
			public.POST("/login", authHandler.Login)
			public.POST("/register", authHandler.Register)
		}

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", func(c *gin.Context) {
				userID := c.GetUint("user_id")
				username := c.GetString("username")
				c.JSON(200, gin.H{
					"user_id":  userID,
					"username": username,
				})
			})
		}

		// 管理员路由组
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		// admin.Use(middleware.AdminAuth()) // 可添加管理员权限中间件
		{
			admin.GET("/users", userHandler.Users)
		}
	}

	return r
}

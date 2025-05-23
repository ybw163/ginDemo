package main

import (
	"errors"
	"gin-web-project/internal/config"
	"gin-web-project/internal/database"
	"gin-web-project/internal/router"
	"gin-web-project/pkg/logger"
	"log"
	"net/http"
	"time"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Path)

	// 连接数据库
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 设置路由
	r := router.Setup(db, cfg)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	logger.Info("Server starting on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}

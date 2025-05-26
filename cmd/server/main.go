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
	config.Load()
	// 初始化日志
	logger.Init()
	// 连接数据库
	database.Connect()
	// 初始化Redis
	database.ConnectRedis()
	//设置路由
	r := router.Setup()
	//创建HTTP服务器
	serverConfig := config.Cfg.Server
	server := &http.Server{
		Addr:         ":" + serverConfig.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second,
	}

	logger.Info("Server starting on port %s", serverConfig.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}

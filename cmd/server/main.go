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
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()

	// 连接数据库
	db, err := database.Connect()
	if err != nil {
		panic("Failed to connect to Mysql: " + err.Error())
	}
	//  初始化Redis
	redisErr := database.ConnectRedis()
	if redisErr != nil {
		panic("Failed to connect to Mysql: " + redisErr.Error())
	}

	// 数据库迁移
	database.InitDB(db)

	// 设置路由
	r := router.Setup(db)

	// 创建HTTP服务器
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

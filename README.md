# 项目目录结构
```
gin-web-project/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go              # 配置结构和加载
│   ├── middleware/
│   │   ├── auth.go                # JWT认证中间件
│   │   ├── cors.go                # 跨域中间件
│   │   ├── logger.go              # 日志中间件
│   │   └── rate_limit.go          # 限流中间件
│   ├── handler/
│   │   ├── auth.go                # 认证处理器
│   │   ├── user.go                # 用户处理器
│   │   └── health.go              # 健康检查
│   ├── service/
│   │   ├── auth.go                # 认证服务
│   │   └── user.go                # 用户服务
│   ├── model/
│   │   ├── user.go                # 用户模型
│   │   └── base.go                # 基础模型
│   ├── database/
│   │   └── database.go            # 数据库连接
│   └── router/
│       └── router.go              # 路由配置
├── pkg/
│   ├── utils/
│   │   ├── jwt.go                 # JWT工具
│   │   ├── response.go            # 响应工具
│   │   └── validator.go           # 验证工具
│   └── logger/
│       └── logger.go              # 日志工具
├── configs/
│   ├── config.yaml                # 配置文件
│   └── config.prod.yaml           # 生产环境配置
├── logs/                          # 日志目录
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 1. go.mod
```go
module gin-web-project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/spf13/viper v1.16.0
    gorm.io/driver/mysql v1.5.1
    gorm.io/gorm v1.25.4
    github.com/go-playground/validator/v10 v10.15.1
    github.com/sirupsen/logrus v1.9.3
    golang.org/x/time v0.3.0
)
```

## 2. cmd/server/main.go - 程序入口
```go
package main

import (
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
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

## 3. internal/config/config.go - 配置管理
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
    Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
    Port         string `mapstructure:"port"`
    Mode         string `mapstructure:"mode"` // debug, release, test
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
    Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime int    `mapstructure:"expire_time"` // hours
}

type LogConfig struct {
    Level string `mapstructure:"level"`
    Path  string `mapstructure:"path"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    viper.AddConfigPath(".")

    // 设置默认值
    setDefaults()

    // 读取环境变量
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func setDefaults() {
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.mode", "debug")
    viper.SetDefault("server.read_timeout", 30)
    viper.SetDefault("server.write_timeout", 30)
    viper.SetDefault("database.driver", "mysql")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", "3306")
    viper.SetDefault("database.charset", "utf8mb4")
    viper.SetDefault("jwt.expire_time", 24)
    viper.SetDefault("log.level", "info")
    viper.SetDefault("log.path", "./logs")
}
```

## 4. internal/database/database.go - 数据库连接
```go
package database

import (
    "fmt"
    "gin-web-project/internal/config"
    "gin-web-project/internal/model"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "time"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
        cfg.Charset,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    // 设置连接池
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return db, nil
}

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        // 添加其他模型
    )
}
```

## 5. internal/model/base.go - 基础模型
```go
package model

import (
    "gorm.io/gorm"
    "time"
)

type BaseModel struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## 6. internal/model/user.go - 用户模型
```go
package model

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
```

## 7. internal/middleware/cors.go - 跨域中间件
```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

## 8. internal/middleware/logger.go - 日志中间件
```go
package middleware

import (
    "gin-web-project/pkg/logger"
    "github.com/gin-gonic/gin"
    "time"
)

func Logger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return logger.Infof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```

## 9. internal/middleware/auth.go - JWT认证中间件
```go
package middleware

import (
    "gin-web-project/pkg/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // 移除 "Bearer " 前缀
        if strings.HasPrefix(token, "Bearer ") {
            token = token[7:]
        }

        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

## 10. internal/middleware/rate_limit.go - 限流中间件
```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
    "net/http"
    "sync"
)

var (
    limiter = rate.NewLimiter(10, 100) // 每秒10个请求，突发100个
    mu      sync.RWMutex
    clients = make(map[string]*rate.Limiter)
)

func RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        mu.RLock()
        limiter, exists := clients[clientIP]
        mu.RUnlock()

        if !exists {
            mu.Lock()
            limiter = rate.NewLimiter(10, 100)
            clients[clientIP] = limiter
            mu.Unlock()
        }

        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 11. pkg/utils/jwt.go - JWT工具
```go
package utils

import (
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "time"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key") // 应该从配置文件读取

func GenerateToken(userID uint, username string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

## 12. pkg/utils/response.go - 响应工具
```go
package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}

func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    400,
        Message: message,
    })
}

func Unauthorized(c *gin.Context, message string) {
    c.JSON(http.StatusUnauthorized, Response{
        Code:    401,
        Message: message,
    })
}

func InternalError(c *gin.Context, message string) {
    c.JSON(http.StatusInternalServerError, Response{
        Code:    500,
        Message: message,
    })
}
```

## 13. pkg/logger/logger.go - 日志工具
```go
package logger

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "os"
    "path/filepath"
)

var log *logrus.Logger

func Init(level, logPath string) {
    log = logrus.New()

    // 设置日志级别
    switch level {
    case "debug":
        log.SetLevel(logrus.DebugLevel)
    case "info":
        log.SetLevel(logrus.InfoLevel)
    case "warn":
        log.SetLevel(logrus.WarnLevel)
    case "error":
        log.SetLevel(logrus.ErrorLevel)
    default:
        log.SetLevel(logrus.InfoLevel)
    }

    // 创建日志目录
    if err := os.MkdirAll(logPath, 0755); err != nil {
        log.Fatalf("Failed to create log directory: %v", err)
    }

    // 设置日志输出文件
    logFile, err := os.OpenFile(filepath.Join(logPath, "app.log"), os.O_CREATE|os.O_WRITEV|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }

    log.SetOutput(logFile)
    log.SetFormatter(&logrus.JSONFormatter{})
}

func Debug(args ...interface{}) {
    log.Debug(args...)
}

func Info(args ...interface{}) {
    log.Info(args...)
}

func Warn(args ...interface{}) {
    log.Warn(args...)
}

func Error(args ...interface{}) {
    log.Error(args...)
}

func Infof(format string, args ...interface{}) string {
    msg := fmt.Sprintf(format, args...)
    log.Info(msg)
    return msg
}
```

## 14. internal/handler/health.go - 健康检查
```go
package handler

import (
    "gin-web-project/pkg/utils"
    "github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
    utils.Success(c, gin.H{
        "status": "ok",
        "message": "Server is running",
    })
}
```

## 15. internal/handler/auth.go - 认证处理器
```go
package handler

import (
    "gin-web-project/internal/model"
    "gin-web-project/internal/service"
    "gin-web-project/pkg/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{
        authService: service.NewAuthService(db),
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req model.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        utils.Error(c, 401, err.Error())
        return
    }

    utils.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req model.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err.Error())
        return
    }

    user, err := h.authService.Register(req.Username, req.Email, req.Password)
    if err != nil {
        utils.Error(c, 400, err.Error())
        return
    }

    utils.Success(c, user)
}
```

## 16. internal/service/auth.go - 认证服务
```go
package service

import (
    "errors"
    "gin-web-project/internal/model"
    "gin-web-project/pkg/utils"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type AuthService struct {
    db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
    return &AuthService{db: db}
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
```

## 17. internal/router/router.go - 路由配置
```go
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
            admin.GET("/users", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "admin users"})
            })
        }
    }

    return r
}
```

## 18. configs/config.yaml - 配置文件
```yaml
server:
  port: "8080"
  mode: "debug"  # debug, release, test
  read_timeout: 30
  write_timeout: 30

database:
  driver: "mysql"
  host: "localhost"
  port: "3306"
  username: "root"
  password: "password"
  database: "gin_web"
  charset: "utf8mb4"

jwt:
  secret: "your-jwt-secret-key"
  expire_time: 24  # hours

log:
  level: "info"  # debug, info, warn, error
  path: "./logs"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

## 19. docker-compose.yml - Docker编排
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
    environment:
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    volumes:
      - ./logs:/app/logs

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: gin_web
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  mysql_data:
```

## 20. Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

CMD ["./main"]
```
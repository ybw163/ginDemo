# GIN的WEB-DEMO用来测试，接入数据库和缓存


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
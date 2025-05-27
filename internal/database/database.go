package database

import (
	"fmt"
	"gin-web-project/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"reflect"
	"time"
)

var (
	RegisteredModels []interface{}
	Db               *gorm.DB
)

func ConnectMySql() {
	cfg := config.Cfg.Database
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
		panic("Failed to connect to Mysql: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to connect to Mysql: " + err.Error())
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	Db = db
	InitDB()
}

func InitDB() {
	for _, model := range RegisteredModels {
		err := Db.AutoMigrate(&model)
		if err != nil {
			log.Println("<UNK>:", err)
			continue
		}
		log.Printf("已自动迁移模型: %v\n", reflect.TypeOf(model).Elem().Name())
	}
}

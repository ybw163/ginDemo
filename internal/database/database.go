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

var RegisteredModels []interface{}

func Connect() (*gorm.DB, error) {
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

func InitDB(db *gorm.DB) {
	for _, model := range RegisteredModels {
		err := db.AutoMigrate(&model)
		if err != nil {
			log.Println("<UNK>:", err)
			continue
		}
		log.Printf("已自动迁移模型: %v\n", reflect.TypeOf(model).Elem().Name())
	}
}

package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"collab-node-platform-backend/config"
	"collab-node-platform-backend/model"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.AppConfig.MysqlDSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境打印SQL
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	DB = db
	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Task{},
		&model.TaskMember{},
		&model.Node{},
		&model.OperationLog{},
	)
	if err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}
	log.Println("数据库连接成功并完成自动迁移")
}

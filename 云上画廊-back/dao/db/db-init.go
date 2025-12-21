package db

import (
	"fmt"
	"os"
	"painting/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// 初始化
func Init() (*Database, bool) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("root:123456@tcp(%s:%s)/paint?charset=utf8mb4&parseTime=True&loc=Local", host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) 
	if err != nil {
		return nil, false
	}
	db.AutoMigrate(&model.User{}, &model.Work{}, &model.Comment{})
	re := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	if re.Ping(re.Context()).Err() != nil {
		return nil, false
	}
	DB := &Database{DB: db, RE: re}
	return DB, true
}

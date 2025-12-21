package db

import (
	"context"
	"fmt"
	"os"
	"painting/model"
	"time"

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
	// MySQL
	host := getenvDefault("DB_HOST", "mysql")
	port := getenvDefault("DB_PORT", "3306")
	user := getenvDefault("DB_USER", "root")
	pass := getenvDefault("DB_PASS", "123456")
	dbname := getenvDefault("DB_NAME", "paint")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		user, pass, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("mysql connect error: %v\n", err)
		return nil, false
	}

	if err := db.AutoMigrate(&model.User{}, &model.Work{}, &model.Comment{}); err != nil {
		fmt.Printf("mysql automigrate error: %v\n", err)
		return nil, false
	}

	// Redis
	rHost := getenvDefault("REDIS_HOST", "redis")
	rPort := getenvDefault("REDIS_PORT", "6379")
	rPass := getenvDefault("REDIS_PASS", "123456")

	re := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rHost, rPort),
		Password: rPass,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := re.Ping(ctx).Err(); err != nil {
		fmt.Printf("redis ping error: %v\n", err)
		return nil, false
	}

	DB := &Database{DB: db, RE: re}
	return DB, true
}

func getenvDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

package db

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
	RE *redis.Client
}

package models

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB = Init()

var RDB = InitRedis()

func Init() *gorm.DB {
	dsn := "root:123456@(39.101.1.158:3003)/gorm_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return db
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "39.101.1.158:3004",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

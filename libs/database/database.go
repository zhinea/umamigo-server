package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zhinea/umamigo-server/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var Redis *redis.Client
var Ctx = context.Background()

func Connect() {
	ConnectMySQL()
	ConnectRedis()
}

func ConnectMySQL() {
	var err error

	DB, err = gorm.Open(mysql.Open(utils.Cfg.Database.MySQL.DSN), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	log.Println("[DATABASE] -> MySQL Connected")
}

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     utils.Cfg.Database.Redis.Addr,
		Password: utils.Cfg.Database.Redis.Password,
		DB:       utils.Cfg.Database.Redis.DB,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}

	log.Println("[DATABASE] -> Redis Connected")
}

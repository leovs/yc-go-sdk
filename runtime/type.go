package runtime

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	redisClient "github.com/leovs/yc-go-sdk/redis-client"
	"gorm.io/gorm"
)

type Runtime interface {
	Mode(value ...string) string
	IsDebug() bool
	SetDb(db *gorm.DB)
	GetDb() *gorm.DB

	SetEs(es *elasticsearch.TypedClient)
	GetEs() *elasticsearch.TypedClient

	// SetEngine 使用的路由
	SetEngine(engine *fiber.App)
	GetEngine() *fiber.App

	GetConfig(key string) interface{}
	SetConfig(key string, value interface{})

	SetRedis(redis *redisClient.RedisClient)
	GetRedis() *redisClient.RedisClient
}

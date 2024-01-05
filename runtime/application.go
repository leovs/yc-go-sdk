package runtime

import (
	"github.com/gofiber/fiber/v2"
	redis_client "github.com/leovs/yc-go-sdk/redis-client"
	"gorm.io/gorm"
	"sync"
)

const debugFlag = "debug"

type Application struct {
	db      *gorm.DB
	redis   *redis_client.RedisClient
	engine  *fiber.App
	mux     sync.RWMutex
	configs map[string]interface{} // 系统参数
	mode    string                 // 运行模式
}

func (e *Application) IsDebug() bool {
	return e.mode == debugFlag
}

func (e *Application) Mode(value ...string) string {
	if len(value) > 0 {
		e.mode = value[0]
	}
	return e.mode
}

// SetDb 设置对应key的db
func (e *Application) SetDb(db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.db = db
}

// GetDb 获取所有map里的db数据
func (e *Application) GetDb() *gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.db
}

// SetEngine 设置路由引擎
func (e *Application) SetEngine(engine *fiber.App) {
	e.engine = engine
}

// GetEngine 获取路由引擎
func (e *Application) GetEngine() *fiber.App {
	return e.engine
}

func (e *Application) SetConfig(key string, value interface{}) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.configs[key] = value
}

func (e *Application) GetConfig(key string) interface{} {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.configs[key]
}

func (e *Application) SetRedis(redis *redis_client.RedisClient) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.redis = redis
}
func (e *Application) GetRedis() *redis_client.RedisClient {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.redis
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		db:      nil,
		configs: make(map[string]interface{}),
	}
}

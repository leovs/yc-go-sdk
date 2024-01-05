// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	_const "github.com/leovs/yc-go-sdk/const"
	redis_client "github.com/leovs/yc-go-sdk/redis-client"
	"github.com/leovs/yc-go-sdk/sdk"
	"gorm.io/gorm"
	"sync/atomic"
)

const cleanCacheScript = `
local keys = redis.call('keys', ARGV[1])
for i=1,#keys,5000 do 
	redis.call('del', 'defaultKey', unpack(keys, i, math.min(i+4999, #keys)))
end
return 1
`

var (
	_ gorm.Plugin = &Gorm2Cache{}
)

type Gorm2CacheConfig struct {
	Enable bool `yaml:"enable"` // 是否开启缓存
}

type Gorm2Cache struct {
	Config     Gorm2CacheConfig
	db         *gorm.DB
	redis      *redis_client.RedisClient
	cleanCache *redis.Script
	hitCount   int64
}

func (g *Gorm2Cache) Name() string {
	return _const.GormCachePrefix
}

func (g *Gorm2Cache) Initialize(db *gorm.DB) (err error) {
	if err = db.Callback().Create().After("*").Register("gorm:cache:after_create", AfterCreate(g)); err != nil {
		return err
	}

	if err = db.Callback().Delete().After("*").Register("gorm:cache:after_delete", AfterDelete(g)); err != nil {
		return err
	}

	if err = db.Callback().Update().After("*").Register("gorm:cache:after_update", AfterUpdate(g)); err != nil {
		return err
	}

	_ = db.Callback().Query().Replace("gorm:query", BeforeQuery(g))

	/*if err = db.Callback().Query().Before("gorm:query").Register("gorm:cache:before_query", BeforeQuery(g)); err != nil {
		return err
	}*/

	if err = db.Callback().Query().After("*").Register("gorm:cache:after_query", AfterQuery(g)); err != nil {
		return err
	}

	g.db = db
	g.redis = sdk.Runtime.GetRedis()
	// 初始化redis脚本
	g.cleanCache = redis.NewScript(1, cleanCacheScript)
	return
}

func (g *Gorm2Cache) GetHitCount() int64 {
	return atomic.LoadInt64(&g.hitCount)
}

func (g *Gorm2Cache) ResetHitCount() {
	atomic.StoreInt64(&g.hitCount, 0)
}

func (g *Gorm2Cache) IncrHitCount() {
	atomic.AddInt64(&g.hitCount, 1)
}

// CleanCache 清理缓存
func (g *Gorm2Cache) CleanCache(key string) error {
	key = fmt.Sprintf("%s:%s*", _const.GormCachePrefix, key)
	_, err := g.cleanCache.Do(g.redis.Pool.Get(), nil, key)
	return err
}

// IsCache 是否需要缓存
func (g *Gorm2Cache) IsCache(db *gorm.DB) bool {
	// 全局开关
	if !g.Config.Enable {
		return false
	}

	if isEnable, ok := db.Get(_const.GormCacheEnablePrefix); !ok || !isEnable.(bool) {
		return false
	}
	return true
}

// TableName 获取表名
func (g *Gorm2Cache) TableName(db *gorm.DB) string {
	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}
	return db.Statement.Table
}

func (g *Gorm2Cache) SetBean(key string, cacheValue any, expires int64) error {
	if err := g.redis.SetEx(key, cacheValue, expires); err != nil {
		return err
	}
	return nil
}

func (g *Gorm2Cache) GetBean(key string, ptr any) error {
	return g.redis.GetObject(key, ptr)
}

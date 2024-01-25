// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package cache

import (
	_const "github.com/leovs/yc-go-sdk/const"
	"github.com/leovs/yc-go-sdk/log"
	"gorm.io/gorm"
)

func AfterQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if !db.DryRun && db.Error == nil && cache.IsCache(db) {
			// 是否已被缓存命中
			if v, ok := db.Get(_const.GormCacheHitPrefix); ok && v.(bool) {
				log.Debug("Gorm2Cache cache hit: %v\n", cache.GetHitCount())
				return
			}

			tableName := cache.TableName(db)
			sql := db.Statement.SQL.String()
			key := GenSearchCacheKey(tableName, sql, db.Statement.Vars...)

			// 获取ttl配置
			var ttl int64 = 0
			if ttlCnf, ok := db.Get(_const.GormCacheTTLPrefix); ok {
				ttl = ttlCnf.(int64)
			}

			// 缓存数据
			if err := cache.SetBean(key, db.Statement.Dest, ttl); err != nil {
				log.Error("Gorm2Cache AfterQuery err: %v\n", err)
				return
			}
		}
	}
}

// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package cache

import (
	_const "github.com/leovs/yc-go-sdk/const"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func BeforeQuery(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error == nil {
			callbacks.BuildQuerySQL(db)

			if !db.DryRun && db.Error == nil {
				tableName := cache.TableName(db)
				sql := db.Statement.SQL.String()
				// 尝试从缓存取数据
				if cache.IsCache(db) {
					key := GenSearchCacheKey(tableName, sql, db.Statement.Vars...)
					if err := cache.GetBean(key, db.Statement.Dest); err == nil {
						cache.IncrHitCount()
						db.Set(_const.GormCacheHitPrefix, true)
						return
					}
				}
				// 未命中缓存，执行查询
				rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, sql, db.Statement.Vars...)
				if err != nil {
					_ = db.AddError(err)
					return
				}
				defer func() {
					_ = db.AddError(rows.Close())
				}()
				gorm.Scan(rows, db, 0)
			}
		}
	}
}

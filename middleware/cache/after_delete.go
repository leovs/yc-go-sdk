// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package cache

import (
	"github.com/leovs/yc-go-sdk/log"
	"gorm.io/gorm"
)

func AfterDelete(cache *Gorm2Cache) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.DryRun || db.Error != nil {
			return
		}
		tableName := cache.TableName(db)
		if err := cache.CleanCache(tableName); err != nil {
			log.Error("Gorm2Cache CleanCache err: %v\n", err)
			_ = db.AddError(err)
			return
		}
		log.Debug("Gorm2Cache AfterDelete: %v\n", tableName)
	}
}

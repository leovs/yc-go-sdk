// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	"github.com/leovs/yc-go-sdk/middleware/cache"
	"github.com/leovs/yc-go-sdk/sdk"
	"github.com/leovs/yc-go-sdk/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

type DatabaseConfig struct {
	Separation bool                   `yaml:"separation"` // 是否开启读写分离
	MaxIdle    int                    `yaml:"maxIdle"`    // 最大空闲连接数
	MaxOpen    int                    `yaml:"maxOpen"`    // 最大打开连接数
	Master     string                 `yaml:"master"`     // 主库
	Slave      []string               `yaml:"slave"`      // 从库
	Cache      cache.Gorm2CacheConfig `yaml:"cache"`      // 缓存配置
}

// Init 初始化配置
func (e *DatabaseConfig) Init(_config *Settings) {
	log.Printf("[%v] 正在连接数据库\n", _config.AppName)

	LogLevel := utils.Ternary(_config.Mode == "debug", logger.Info, logger.Warn).(logger.LogLevel)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      LogLevel,    // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: e.Master,
	}), &gorm.Config{Logger: newLogger, PrepareStmt: true, SkipDefaultTransaction: false})
	if err != nil {
		log.Panicf("数据库连接失败 %v\n", err.Error())
		return
	}
	var replicas []gorm.Dialector
	for i, s := range e.Slave {
		log.Printf("读写分离-%d-%s \n", i, s)
		replicas = append(replicas, mysql.New(mysql.Config{DSN: s}))
	}

	err = db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: e.Master,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(e.MaxIdle).
			SetMaxOpenConns(e.MaxOpen),
	)

	if err != nil {
		log.Panicf("数据库连接失败 %v\n", err.Error())
		return
	}

	// 配置缓存
	if e.Cache.Enable {
		err := db.Use(&cache.Gorm2Cache{Config: e.Cache})
		if err != nil {
			log.Panicf("数据库缓存初始化失败 %v\n", err.Error())
			return
		}
		log.Println("数据库缓存初始化成功")
	}

	sdk.Runtime.SetDb(db)
}

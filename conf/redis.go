// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	"github.com/leovs/yc-go-sdk/log"
	redis_client "github.com/leovs/yc-go-sdk/redis-client"
	"github.com/leovs/yc-go-sdk/sdk"
)

type RedisConfig struct {
	MasterName     string `yaml:"masterName"`     // Sentinel 哨兵模式 Master名字
	Address        string `yaml:"address"`        // 地址 localhost:6379
	Password       string `yaml:"password"`       // 密码
	DBIds          int    `yaml:"dbIds"`          // redisDB
	MaxIdle        int    `yaml:"maxIdle"`        // redis连接池最大空闲连接数
	MaxActive      int    `yaml:"maxActive"`      // redis连接池最大激活连接数, 0为不限制
	ConnectTimeout int    `yaml:"connectTimeout"` // redis连接超时时间, 单位毫秒
	ReadTimeout    int    `yaml:"readTimeout"`    // redis读取超时时间, 单位毫秒
	WriteTimeout   int    `yaml:"writeTimeout"`   // redis写入超时时间, 单位毫秒
}

// Init 初始化配置
func (e *RedisConfig) Init(config *Settings) {
	log.Info("[%v] 正在连接Redis", config.AppName)
	redis := &redis_client.RedisClient{}
	redis.InitRedis(
		e.MasterName,
		e.Address,
		e.Password,
		e.DBIds,
		e.MaxIdle,
		e.MaxActive,
		e.ConnectTimeout,
		e.ReadTimeout,
		e.WriteTimeout)
	sdk.Runtime.SetRedis(redis)
}

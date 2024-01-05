// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	eureka_client "github.com/leovs/yc-go-sdk/eureka-client"
	"github.com/leovs/yc-go-sdk/log"
)

type EurekaConfig struct {
	DefaultZone string `yaml:"defaultZone"`
}

// Init 初始化配置
func (e *EurekaConfig) Init(config *Settings) {
	log.Info("[%v] 正在注册Eureka", config.AppName)
	eureka_client.GetInstance().RunEureka(e.DefaultZone, config.AppName, config.Version, config.Port)
}

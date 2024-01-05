// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leovs/yc-go-sdk/log"
	"github.com/leovs/yc-go-sdk/middleware"
	"github.com/leovs/yc-go-sdk/middleware/validators"
	"github.com/leovs/yc-go-sdk/sdk"
)

type ServerConfig struct {
}

// Init 初始化配置
func (e *ServerConfig) Init(config *Settings) {
	// 初始化程序并设置运行模式
	log.Info("[%v] 版本[%v] 运行模式[%v]", config.AppName, config.Version, config.Mode)
	log.DebugMode = config.Mode == "debug"
	sdk.Runtime.IsDebug()
	engine := fiber.New()
	sdk.Runtime.SetEngine(engine)
	validators.Init()
	middleware.InitMiddleware(engine)
}

// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/leovs/yc-go-sdk/sdk"
	"log"
	"os"
)

type EsConfig struct {
	Addresses          []string `yaml:"addresses"`          // 地址
	Username           string   `yaml:"username"`           // 用户名
	Password           string   `yaml:"password"`           // 密码
	EnableDebugLogger  bool     `yaml:"enableDebugLogger"`  // 是否开启调试日志
	EnableRequestBody  bool     `yaml:"enableRequestBody"`  // 是否开启调试日志
	EnableResponseBody bool     `yaml:"enableResponseBody"` // 是否开启调试日志
}

// Init 初始化配置
func (e *EsConfig) Init(_config *Settings) {
	log.Printf("[%v] 正在连接ES\n", _config.AppName)
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses:         _config.EsConfig.Addresses,
		Username:          _config.EsConfig.Username,
		Password:          _config.EsConfig.Password,
		EnableDebugLogger: _config.EsConfig.EnableDebugLogger,
		Logger: &elastictransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  _config.EsConfig.EnableRequestBody,
			EnableResponseBody: _config.EsConfig.EnableResponseBody,
		},
	})
	if err != nil {
		log.Panicf("ES连接失败 %v\n", err.Error())
	}
	sdk.Runtime.SetEs(client)
}

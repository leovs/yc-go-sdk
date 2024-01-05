// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package conf

import (
	"fmt"
	"github.com/leovs/yc-go-sdk/log"
	"github.com/leovs/yc-go-sdk/sdk"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"time"
)

const (
	cfgFile  = "config.%s.yaml"
	LOCATION = "Asia/Shanghai"
)

type Settings struct {
	Mode           string          `yaml:"mode"`       // 运行模式
	AppName        string          `yaml:"appName"`    // 应用名称
	Port           int             `yaml:"port"`       // 端口
	Version        string          `yaml:"version"`    // 版本
	EurekaConfig   *EurekaConfig   `yaml:"eureka"`     // 服务注册
	RedisConfig    *RedisConfig    `yaml:"redis"`      // Redis配置
	ServerConfig   *ServerConfig   `yaml:"server"`     // Gin配置
	DatabaseConfig *DatabaseConfig `yaml:"dataSource"` // 数据库配置
}

// Setup 载入配置文件
func (e *Settings) Setup(env string) {
	conStr, err := os.ReadFile(fmt.Sprintf(cfgFile, env))
	if err != nil {
		log.Panic("读取配置文件失败", err)
		return
	}

	if err := yaml.Unmarshal(conStr, &e); err != nil {
		log.Panic("解析配置文件失败", err)
		return
	}

	// 初始化配置
	e.Init()
}

// 初始化全局变量
func (e *Settings) initConfig() {
	configElem := reflect.ValueOf(e).Elem()
	relType := configElem.Type()
	for i := 0; i < relType.NumField(); i++ {
		name := relType.Field(i).Name
		sdk.Runtime.SetConfig(name, configElem.Field(i).Interface())
	}
	sdk.Runtime.Mode(e.Mode)
}

// Init 初始化配置
func (e *Settings) Init() {
	// 设置时区
	_, _ = time.LoadLocation(LOCATION)

	// 初始化参数
	e.initConfig()
	// 初始化Redis
	e.RedisConfig.Init(e)
	// 初始化数据库
	e.DatabaseConfig.Init(e)
	// eureka注册
	e.EurekaConfig.Init(e)
	// service注册
	e.ServerConfig.Init(e)
}

func (e *Settings) SetRouter(routers ...sdk.IRouter) {
	engine := sdk.Runtime.GetEngine()
	for _, router := range routers {
		router.InitRouter(&sdk.Router{Engine: engine})
	}
}

func (e *Settings) Run() {
	err := sdk.Runtime.GetEngine().Listen(fmt.Sprintf(":%d", e.Port))
	if err != nil {
		panic(err)
	}
}

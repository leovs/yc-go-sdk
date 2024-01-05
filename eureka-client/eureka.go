// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package eureka_client

import (
	"github.com/go-resty/resty/v2"
	"sync"
)

type Eureka struct {
	eureka *Client
	client *resty.Client
}

type Url struct {
	Protocol   string
	ServerName string
	ServerPort int
	Path       string
	Param      string
}

var (
	instance *Eureka
	once     sync.Once
)

func GetInstance() *Eureka {
	once.Do(func() {
		instance = &Eureka{}
	})
	return instance
}

// RunEureka 运行eureka
func (c *Eureka) RunEureka(defaultZone, app, version string, port int) {
	c.eureka = NewClient(&Config{
		DefaultZone:           defaultZone,
		App:                   app,
		Port:                  port,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              version,
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	}).Start()
}

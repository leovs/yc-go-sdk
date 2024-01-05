// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package _const

import "github.com/leovs/yc-go-sdk/errors"

type EnumCommon struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

const (
	GormCachePrefix       = "gorm2cache"
	GormCacheEnablePrefix = "Gorm2Cache:enable"
	GormCacheTTLPrefix    = "Gorm2Cache:ttl"
	GormCacheHitPrefix    = "Gorm2Cache:hit"
)

// 服务接口错误返回
var (
	Failure        = &errors.Message{Code: -1, Msg: "失败", Data: nil}
	Success        = &errors.Message{Code: 0, Msg: "成功", Data: nil}
	ParamError     = &errors.Message{Code: -2, Msg: "参数错误", Data: nil}
	SystemError    = &errors.Message{Code: -3, Msg: "系统异常", Data: nil}
	NoDataReturn   = &errors.Message{Code: -4, Msg: "暂无数据", Data: nil}
	DbConnectError = &errors.Message{Code: -5, Msg: "数据库连接获取失败", Data: nil}
	DataExisted    = &errors.Message{Code: -1003, Msg: "数据已存在", Data: nil}
	DataReferenced = &errors.Message{Code: -4001, Msg: "关联数据不存在，无法进行此操作", Data: nil}
)

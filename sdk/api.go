// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package sdk

import (
	_const "github.com/leovs/yc-go-sdk/const"
	"github.com/leovs/yc-go-sdk/errors"
	"github.com/leovs/yc-go-sdk/utils"
)

type Api struct {
}

type ApiMessage errors.Message

// Success 通常成功数据处理
func (e *Api) Success(data ...any) *errors.Message {
	return _const.Success.SetData(utils.Ternary(len(data) == 1, data[0], data))
}

// Failure 通常错误数据处理
func (e *Api) Failure(data ...any) *errors.Message {
	return _const.Failure.SetData(utils.Ternary(len(data) == 1, data[0], data))
}

func (e *Api) Errors(err error) *errors.Message {
	if msg, ok := err.(*errors.Message); ok {
		return msg
	}
	return _const.SystemError.SetMsg(err.Error())
}

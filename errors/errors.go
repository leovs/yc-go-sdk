// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package errors

import "fmt"

type Message struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 状态消息
	Data interface{} `json:"data"`
}

func (e *Message) Error() string {
	return e.Msg
}

func (e *Message) ErrorCode() int {
	return e.Code
}

func (e *Message) SetErrorCode(code int) *Message {
	return &Message{Code: code, Data: e.Data, Msg: e.Msg}
}
func (e *Message) SetMsg(s string, prs ...interface{}) *Message {
	if len(prs) == 0 {
		return &Message{Code: e.Code, Data: e.Data, Msg: s}
	}
	return &Message{Code: e.Code, Data: e.Data, Msg: fmt.Sprintf(s, prs)}
}

func (e *Message) SetData(d interface{}) *Message {
	return &Message{Code: e.Code, Data: d, Msg: e.Msg}
}

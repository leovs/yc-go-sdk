// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package sdk

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_const "github.com/leovs/yc-go-sdk/const"
	"github.com/leovs/yc-go-sdk/errors"
	"runtime/debug"

	//"github.com/leovs/yc-go-sdk/errors"
	"github.com/leovs/yc-go-sdk/middleware/validators"
	"github.com/leovs/yc-go-sdk/runtime"
)

const (
	MethodGet   = "GET"
	MethodPost  = "POST"
	WithContext = "WithContext"
)

type IRouter interface {
	InitRouter(r *Router)
}

type Router struct {
	Engine *fiber.App
	Router fiber.Router
}

func (r *Router) Group(relativePath string) *Router {
	r.Router = r.Engine.Group(relativePath)
	return r
}

func (r *Router) Exec(c *fiber.Ctx, method string, handlers any) error {
	reflectFunc := runtime.ReflectFunc{}
	models := reflectFunc.Of(handlers)
	if len(models) != 2 {
		fmt.Printf("注入对象失败, 参数约定model,service")
		return _const.Failure
	}

	var Model = models[0]
	var Service = models[1]

	// 解析参数
	if method == MethodGet {
		_ = c.QueryParser(Model)
	} else if method == MethodPost {
		_ = c.BodyParser(Model)
	}

	// 验证参数
	if err := validators.Check(Model); err != nil {
		return c.JSON(_const.ParamError.SetData(err))
	}

	// 初始化service
	reflectBind := runtime.ReflectFunc{}
	reflectBind.Of(Service)
	reflectBind.CallMethod(WithContext, c)

	defer func() {
		// 释放对象
		Model = nil
		Service = nil
		models = nil
		// 捕获异常
		if err := recover(); err != nil {
			_ = c.JSON(&errors.Message{Code: -1, Msg: fmt.Sprintf("%+v", err), Data: nil})
			fmt.Printf("panic error=%v, stack=%s \n", err, debug.Stack())
		}
	}()

	result := reflectFunc.Call(Model, Service)
	if len(result) == 1 {
		return c.JSON(result[0].Interface())
	}
	return c.JSON(_const.SystemError)
}

func (r *Router) GET(relativePath string, handlers any) {
	r.Router.Get(relativePath, func(c *fiber.Ctx) error {
		return r.Exec(c, MethodGet, handlers)
	})
}

func (r *Router) POST(relativePath string, handlers any) {
	r.Router.Post(relativePath, func(c *fiber.Ctx) error {
		return r.Exec(c, MethodPost, handlers)
	})
}

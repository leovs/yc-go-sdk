package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leovs/yc-go-sdk/sdk"
)

// WithContextDb 数据库链接
func WithContextDb(c *fiber.Ctx) error {
	if db := sdk.Runtime.GetDb(); db != nil {
		c.Locals("db", db.WithContext(c.Context()))
	}
	return c.Next()
}

// InitMiddleware 初始化中间件
func InitMiddleware(r *fiber.App) {
	// 数据库链接
	r.Use(WithContextDb)
}

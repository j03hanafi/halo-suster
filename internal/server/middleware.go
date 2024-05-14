package server

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/id"
	"github.com/j03hanafi/halo-suster/common/logger"
)

const (
	requestId   = "requestId"
	accessToken = "accessToken"
)

func setMiddlewares(app *fiber.App) {
	app.Use(compressionMiddleware())
	app.Use(recoveryMiddleware())
	app.Use(zapMiddleware())
	app.Use(requestIDMiddleware())
	app.Use(loggerMiddleware())
}

func compressionMiddleware() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	})
}

func recoveryMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

func zapMiddleware() fiber.Handler {
	return fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
		Fields: []string{
			"latency",
			"time",
			"requestId",
			"pid",
			"status",
			"method",
			"path",
			"queryParams",
			"body",
			"ip",
			"ua",
			"resBody",
			"error",
		},
	})
}

func requestIDMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Generator: func() string {
			return id.New().String()
		},
		ContextKey: requestId,
	})
}

func loggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		l := zap.L().With(zap.String(requestId, c.Locals(requestId).(string)))
		ctx = logger.WithCtx(ctx, l)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

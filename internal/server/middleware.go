package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/common/id"
	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/common/security"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

const (
	requestId   = "requestId"
	accessToken = "accessToken"
)

func setMiddlewares(app *fiber.App) {
	app.Use(recoveryMiddleware())
	app.Use(zapMiddleware())
	app.Use(requestIDMiddleware())
	app.Use(loggerMiddleware())
	app.Use(pprofMiddleware())
}

func pprofMiddleware() fiber.Handler {
	return pprof.New()
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
		SkipBody: func(c *fiber.Ctx) bool {
			return strings.HasSuffix(c.Path(), "image")
		},
		Levels: []zapcore.Level{zapcore.ErrorLevel, zapcore.ErrorLevel, zapcore.InfoLevel},
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

func jwtMiddleware(jwtCache *cache.Cache) fiber.Handler {
	return jwtware.New(jwtware.Config{
		Filter: func(c *fiber.Ctx) bool {
			token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
			if token == "" {
				return false
			}
			user, found := jwtCache.Get(token)
			if !found {
				return false
			}

			c.Locals(domain.UserFromToken, user.(domain.User))
			return true
		},
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(configs.Get().JWT.JWTSecret),
		},
		Claims:     &security.AccessTokenClaims{},
		ContextKey: accessToken,
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals(accessToken).(*jwt.Token)
			claims := token.Claims.(*security.AccessTokenClaims)
			user := domain.User{
				ID:   claims.User.UserID,
				NIP:  claims.User.NIP,
				Name: claims.User.Name,
				Role: claims.User.Role,
			}

			exp, _ := claims.GetExpirationTime()
			cacheExp := time.Until(exp.Time) - time.Minute
			if cacheExp > 0 {
				jwtCache.Set(token.Raw, user, cacheExp)
			}

			c.Locals(domain.UserFromToken, user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})
}

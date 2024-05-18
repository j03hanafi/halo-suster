package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/adapter"
	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application"
)

func Run() {
	callerInfo := "[server.Run]"
	l := zap.L().With(zap.String("caller", callerInfo))

	db := adapter.GetDBPool()
	defer db.Close()

	s3 := adapter.GetS3Client()

	jwtCache := adapter.GetJWTCache()
	defer jwtCache.Flush()

	serverTimeout := time.Duration(configs.Get().API.Timeout) * time.Second
	serverConfig := fiber.Config{
		AppName:                   configs.Get().App.Name,
		DisableDefaultDate:        true,
		EnablePrintRoutes:         true,
		JSONDecoder:               json.Unmarshal,
		JSONEncoder:               json.Marshal,
		ReadTimeout:               serverTimeout,
		CaseSensitive:             true,
		StrictRouting:             true,
		DisableHeaderNormalizing:  true,
		DisableDefaultContentType: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := http.StatusInternalServerError

			var fiberErr *fiber.Error
			var handlerErr ErrorHandler

			switch {
			case errors.As(err, &handlerErr):
				code = handlerErr.Status()
			case errors.As(err, &fiberErr):
				code = fiberErr.Code
			}

			if code == http.StatusRequestEntityTooLarge {
				return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
					"message": "Request body too large",
				})
			}

			return ctx.Status(code).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	}

	if configs.Get().App.PreFork {
		serverConfig.Prefork = true
	}

	app := fiber.New(serverConfig)
	setMiddlewares(app)
	application.New(app, db, s3, jwtCache, jwtMiddleware(jwtCache))
	l.Debug("Server Config", zap.Any("Config", app.Config()))

	go func() {
		addr := fmt.Sprintf("%s:%d", configs.Get().App.Host, configs.Get().App.Port)
		if err := app.Listen(addr); err != nil {
			l.Panic("Server Error", zap.Error(err))
		}
	}()

	l.Info("Server is starting...")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	if err := app.ShutdownWithTimeout(serverTimeout); err != nil {
		l.Panic("Server forced to shutdown", zap.Error(err))
	}

	l.Info("Server was successful shutdown")
}

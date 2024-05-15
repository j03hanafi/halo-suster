package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/common/database"
	"github.com/j03hanafi/halo-suster/internal/application"
)

func Run() {
	callerInfo := "[server.Run]"
	l := zap.L().With(zap.String("caller", callerInfo))

	db, err := database.NewPGConn()
	if err != nil {
		l.Panic("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	serverTimeout := time.Duration(configs.Get().API.Timeout) * time.Second
	serverConfig := fiber.Config{
		AppName:                   configs.Get().App.Name,
		DisableDefaultDate:        true,
		EnablePrintRoutes:         true,
		JSONDecoder:               sonic.Unmarshal,
		JSONEncoder:               sonic.Marshal,
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
	application.New(app, db, jwtMiddleware())
	l.Debug("Server Config", zap.Any("Config", app.Config()))

	go func() {
		addr := fmt.Sprintf("%s:%d", configs.Get().App.Host, configs.Get().App.Port)
		if err = app.Listen(addr); err != nil {
			l.Panic("Server Error", zap.Error(err))
		}
	}()

	l.Info("Server is starting...")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	err = app.ShutdownWithTimeout(serverTimeout)
	if err != nil {
		l.Panic("Server forced to shutdown", zap.Error(err))
	}

	l.Info("Server was successful shutdown")
}

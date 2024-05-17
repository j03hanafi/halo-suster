package internal

import (
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/server"
)

func Run() {
	callerInfo := "[internal.Run]"

	l := logger.Get()
	defer func() {
		_ = l.Sync()
	}()
	zap.ReplaceGlobals(l)

	l.With(zap.String("caller", callerInfo)).Debug("config loaded", zap.Any("config", configs.Get()))

	server.Run()
}

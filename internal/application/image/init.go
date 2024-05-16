package image

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/image/handler"
	"github.com/j03hanafi/halo-suster/internal/application/image/service"
)

func NewModule(router fiber.Router, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Get().App.ContextTimeout) * time.Second

	imageService := service.NewImageService(ctxTimeout)
	handler.NewImageHandler(router, jwtMiddleware, imageService)
}

package image

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/image/handler"
	"github.com/j03hanafi/halo-suster/internal/application/image/repository"
	"github.com/j03hanafi/halo-suster/internal/application/image/service"
)

func NewModule(router fiber.Router, s3 *s3.Client, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Get().App.ContextTimeout) * time.Second

	imageRepository := repository.NewImageRepository(s3)
	imageService := service.NewImageService(ctxTimeout, imageRepository)
	handler.NewImageHandler(router, jwtMiddleware, imageService)
}

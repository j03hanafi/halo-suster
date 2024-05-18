package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/application/image/service"
)

type imageHandler struct {
	imageService service.ImageServiceContract
}

func NewImageHandler(router fiber.Router, jwtMiddleware fiber.Handler, imageService service.ImageServiceContract) {
	handler := imageHandler{
		imageService: imageService,
	}

	imageRouter := router.Group("/image", jwtMiddleware)
	imageRouter.Post("", handler.UploadImage)
}

func (h imageHandler) UploadImage(c *fiber.Ctx) error {
	callerInfo := "[imageHandler.UpdateImage]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	l.Debug("body", zap.Any("body", c.Body()))

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	// Validate File
	const minSize, maxSize = 10 * 1024, 2 * 1024 * 1024
	if file.Size < int64(minSize) || file.Size > int64(maxSize) {
		l.Error("invalid file size")
		return errBadRequest{err: errors.New("file size must be between 10KB and 2MB")}
	}

	if file.Header.Get("Content-Type") != "image/jpeg" {
		l.Error("invalid file type")
		return errBadRequest{err: errors.New("file type must be image/jpeg")}
	}

	url, err := h.imageService.UploadImage(userCtx, file)
	if err != nil {
		l.Error("error uploading image", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "File uploaded successfully"

	imgURL := imageUploadResAcquire()
	defer imageUploadResRelease(imgURL)

	imgURL.ImgURL = url

	res.Data = imgURL

	return c.JSON(res)
}

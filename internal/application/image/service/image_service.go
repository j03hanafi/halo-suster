package service

import (
	"context"
	"mime/multipart"
	"time"

	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
)

type ImageService struct {
	contextTimeout time.Duration
}

func NewImageService(timeout time.Duration) *ImageService {
	return &ImageService{
		contextTimeout: timeout,
	}
}

func (s ImageService) UploadImage(ctx context.Context, image *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[ImageService.UploadImage]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	l.Info("uploading image", zap.String("filename", image.Filename))

	return nil
}

var _ ImageServiceContract = (*ImageService)(nil)

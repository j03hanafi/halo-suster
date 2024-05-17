package service

import (
	"context"
	"mime/multipart"
)

type ImageServiceContract interface {
	UploadImage(ctx context.Context, image *multipart.FileHeader) (string, error)
}

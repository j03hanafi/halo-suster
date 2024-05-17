package repository

import (
	"context"
	"mime/multipart"
)

type ImageRepositoryContract interface {
	UploadImage(ctx context.Context, image *multipart.FileHeader) (string, error)
}

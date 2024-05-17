package adapter

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/configs"
)

var (
	s3Client *s3.Client
	s3Once   sync.Once
)

func GetS3Client() *s3.Client {
	s3Once.Do(func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(configs.Get().App.ContextTimeout)*time.Second,
		)
		defer cancel()

		callerInfo := "[adapter.GetS3Client]"
		l := zap.L().With(zap.String("caller", callerInfo))

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			l.Error("error loading config", zap.Error(err))
			panic(err)
		}

		s3Client = s3.NewFromConfig(cfg)
	})

	return s3Client
}

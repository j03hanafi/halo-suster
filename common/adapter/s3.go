package adapter

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
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

		cfg, err := config.LoadDefaultConfig(ctx, config.WithLogger(&s3log{logger: zap.L()}))
		if err != nil {
			l.Error("error loading config", zap.Error(err))
			panic(err)
		}

		s3Client = s3.NewFromConfig(cfg)
	})

	return s3Client
}

type s3log struct {
	logger *zap.Logger
}

func (s *s3log) Logf(classification logging.Classification, format string, v ...interface{}) {
	switch classification {
	case logging.Warn:
		s.logger.Warn(format, zap.Any("message", v))
	case logging.Debug:
		s.logger.Debug(format, zap.Any("message", v))
	}
}

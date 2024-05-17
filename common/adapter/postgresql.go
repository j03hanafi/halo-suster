package adapter

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	pgxZap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/configs"
)

var (
	pgxPool *pgxpool.Pool
	pgxOnce sync.Once
)

func GetDBPool() *pgxpool.Pool {
	pgxOnce.Do(func() {
		var err error
		pgxPool, err = newPGConn()
		if err != nil {
			panic(err)
		}
	})
	return pgxPool
}

func newPGConn() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(configs.Get().App.ContextTimeout)*time.Second,
	)
	defer cancel()

	callerInfo := "[adapter.newPGConn]"
	l := zap.L().With(zap.String("caller", callerInfo))

	maxConnPool := configs.Get().DB.MaxConnPool
	if configs.Get().DB.MaxConnPoolPercent != 0 && configs.Get().DB.MaxConnPoolPercent <= 1 {
		var err error
		maxConnPool, err = getMaxConnPool(ctx)
		if err != nil {
			l.Error("error getting max connection pool", zap.Error(err))
			return nil, err
		}
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		configs.Get().DB.Username,
		configs.Get().DB.Password,
		configs.Get().DB.Host,
		configs.Get().DB.Port,
		configs.Get().DB.Name,
		configs.Get().DB.Params,
	)

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		l.Error("error parsing database config",
			zap.Error(err),
		)
		return nil, err
	}
	config.MaxConns = int32(maxConnPool)

	poolLog := pgxZap.NewLogger(zap.L())
	poolTracer := &tracelog.TraceLog{
		Logger:   poolLog,
		LogLevel: tracelog.LogLevelNone,
	}
	if configs.Get().App.DebugMode {
		poolTracer.LogLevel = tracelog.LogLevelDebug
	}
	config.ConnConfig.Tracer = poolTracer

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		l.Error("error creating database pool",
			zap.Error(err),
		)
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		l.Error("error pinging database",
			zap.Error(err),
		)
		return nil, err
	} else {
		l.Info("connected to database")
	}

	l.Debug("Database Config", zap.Any("Config", pool.Config().ConnString()))
	return pool, nil
}

func getMaxConnPool(ctx context.Context) (int, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		configs.Get().DB.Username,
		configs.Get().DB.Password,
		configs.Get().DB.Host,
		configs.Get().DB.Port,
		configs.Get().DB.Name,
		configs.Get().DB.Params,
	)

	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return 0, err
	}
	defer dbPool.Close()

	var getMaxConn string
	query := `SHOW max_connections`
	err = dbPool.QueryRow(ctx, query).Scan(&getMaxConn)
	if err != nil {
		return 0, err
	}

	maxConn, err := strconv.ParseFloat(getMaxConn, 64)
	if err != nil {
		return 0, err
	}
	maxConnPool := int(math.Floor(maxConn * configs.Get().DB.MaxConnPoolPercent))
	if maxConnPool < 1 {
		maxConnPool = 1
	}

	return maxConnPool, nil
}

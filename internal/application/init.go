package application

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/image"
	"github.com/j03hanafi/halo-suster/internal/application/info"
	"github.com/j03hanafi/halo-suster/internal/application/medical"
	"github.com/j03hanafi/halo-suster/internal/application/user"
)

func New(server *fiber.App, db *pgxpool.Pool, s3 *s3.Client, jwtCache *cache.Cache, jwtMiddleware fiber.Handler) {
	router := server.Group(configs.Get().API.BaseURL)

	info.NewModule(router, db)
	user.NewModule(router, db, jwtCache, jwtMiddleware)
	medical.NewModule(router, db, jwtMiddleware)
	image.NewModule(router, s3, jwtMiddleware)
}

package security

import (
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/j03hanafi/halo-suster/common/configs"
)

func HashPassword(password string) (string, error) {
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), configs.Get().API.BCryptSalt)
	return utils.UnsafeString(bhash), err
}

func ComparePassword(storedPassword, suppliedPassword string) error {
	return bcrypt.CompareHashAndPassword(utils.UnsafeBytes(storedPassword), utils.UnsafeBytes(suppliedPassword))
}

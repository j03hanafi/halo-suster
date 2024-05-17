package security

import (
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type AccessTokenClaims struct {
	User user `json:"user,omitempty"`
	jwt.RegisteredClaims
}

type user struct {
	UserID ulid.ULID `json:"user_id"`
	NIP    string    `json:"nip"`
	Name   string    `json:"name"`
	Role   string    `json:"role"`
}

func GenerateAccessToken(u *domain.User) (string, error) {
	callerInfo := "[security.GenerateAccessToken]"
	l := zap.L().With(zap.String("caller", callerInfo))

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(configs.Get().JWT.Expire) * time.Second)

	claims := AccessTokenClaims{
		User: user{
			UserID: u.ID,
			NIP:    u.NIP,
			Name:   u.Name,
			Role:   u.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			NotBefore: jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(utils.UnsafeBytes(configs.Get().JWT.JWTSecret))
	if err != nil {
		l.Error("failed to sign token", zap.Error(err))
		return "", err
	}

	return signedString, nil
}

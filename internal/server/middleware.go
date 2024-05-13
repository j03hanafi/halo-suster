package server

import "github.com/gofiber/fiber/v2"

func jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

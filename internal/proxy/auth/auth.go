package auth

import (
	"log"

	"github.com/capcom6/swarm-gateway/internal/common"
	"github.com/capcom6/swarm-gateway/internal/proxy/auth/basic"
	"github.com/capcom6/swarm-gateway/internal/proxy/resolver"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
}

func New(config ...Config) fiber.Handler {

	return func(c *fiber.Ctx) error {
		service := c.Locals(resolver.LocalsKeyService).(common.Service)

		return auth(service.Auth.Type, service.Auth.Data)(c)
	}

}

func auth(authType, data string) fiber.Handler {
	if authType == "" {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if authType == "basic" {
		return basic.New(data)
	}

	log.Printf("unknown auth type: %s", authType)
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

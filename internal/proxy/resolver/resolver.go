package resolver

import (
	"errors"

	"github.com/capcom6/swarm-gateway/internal/repository"
	"github.com/gofiber/fiber/v2"
)

const (
	LocalsKeyService = "service"
)

type Config struct {
}

func New(services *repository.ServicesRepository, config ...Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := c.Get("Host")
		if host == "" {
			return fiber.ErrBadRequest
		}

		service, err := services.GetServiceByHost(host)
		if errors.Is(err, repository.ErrSeviceNotFound) {
			return fiber.ErrBadGateway
		}

		c.Locals("service", service)

		return c.Next()
	}
}

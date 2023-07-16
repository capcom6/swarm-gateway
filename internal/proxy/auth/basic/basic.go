package basic

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/tg123/go-htpasswd"
)

func New(options string) fiber.Handler {
	auth, err := htpasswd.NewFromReader(strings.NewReader(options), htpasswd.DefaultSystems, nil)
	// users := options["users"]

	return basicauth.New(basicauth.Config{
		Authorizer: func(username string, password string) bool {
			if err != nil {
				return false
			}

			return auth.Match(username, password)
		},
	})
}

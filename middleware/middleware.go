package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func Cache(duration time.Duration) fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   duration,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Include query parameters in cache key
			return c.Path() + "?" + string(c.Request().URI().QueryString())
		},
	})
}

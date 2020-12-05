// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// Special thanks to : https://github.com/LeafyCode/express-firebase-auth

package gofiberfirebaseauth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// New - Main
func New(config Config) fiber.Handler {
	// Init config
	cfg := configDefault(config)
	// Return authed handler
	return func(c *fiber.Ctx) error {

		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		// 1) Resolve current route
		url := c.Method() + " " + c.Path()
		// TODO add support for route with params and quarries
		// r := c.Route()
		// fmt.Println(r.Method, r.Path, r.Params, r.Handlers)

		// 2) Compare with current route
		if cfg.IgnoreUrls != nil && len(cfg.IgnoreUrls) > 0 {
			for i := range cfg.IgnoreUrls {
				if cfg.IgnoreUrls[i] == url {
					return c.Next()
				}
			}
		}

		// 3) get token from header
		IDToken := c.Get(fiber.HeaderAuthorization)

		if len(IDToken) == 0 {
			return cfg.ErrorHandler(c, errors.New("Missing Token"))
		}

		// 4) Validate the IdToken
		IsPass, err := cfg.Authorizer(IDToken)

		// 5) IF Id token passed
		if IsPass {
			return c.Next()
		}

		return cfg.ErrorHandler(c, err)
	}
}

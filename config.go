package gofiberfirebaseauth

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for goFiber middleware
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Token authorizer
	Authorizer func(string, string) (bool, error)

	// Skip Email Check.
	// Optional. Default: nil
	CheckEmailVerified bool

	// New firebase authntication object
	// Mandatory. Default: nil
	FirebaseApp *firebase.App

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler fiber.ErrorHandler

	// Ignore urls array
	IgnoreUrls []string

	// Ignore email verification for these routes
	CheckEmailVerifiedIgnoredUrls []string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:                          nil,
	IgnoreUrls:                    nil,
	Authorizer:                    nil,
	ErrorHandler:                  nil,
	SuccessHandler:                nil,
	CheckEmailVerified:            false,
	CheckEmailVerifiedIgnoredUrls: nil,
}

// Initializer
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.FirebaseApp == nil {
		fmt.Println("****************************************************************")
		fmt.Println("gofiberfirebaseauth :: Error PLEASE PASS Firebase App in Config")
		fmt.Println("*****************************************************************")
	}

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Default Authorizer function
	if cfg.Authorizer == nil {
		cfg.Authorizer = func(IDToken string, CurrentURL string) (bool, error) {
			if cfg.FirebaseApp == nil {
				return false, errors.New("Missing Firebase App Object")
			}
			client, err := cfg.FirebaseApp.Auth(context.Background())
			// Verify IDToken
			token, err := client.VerifyIDToken(context.Background(), IDToken)

			fmt.Println(CurrentURL)
			// Throw error for bad token
			if err != nil {
				return false, errors.New("Malformed Token")
			}

			// IF CheckEmailVerified enable in config check email is verified
			if cfg.CheckEmailVerified {
				checkEmail := false
				if cfg.CheckEmailVerifiedIgnoredUrls != nil && len(cfg.CheckEmailVerifiedIgnoredUrls) > 0 {
					for i := range cfg.IgnoreUrls {
						if cfg.CheckEmailVerifiedIgnoredUrls[i] == CurrentURL {
							checkEmail = true
						}
					}
				}
				fmt.Println(token.Claims["email_verified"].(bool))
				if checkEmail {
					// Claim email_verified from token
					if !token.Claims["email_verified"].(bool) {
						return false, errors.New("Email not verified")
					}
				}
			}

			return true, nil
		}
	}

	// Default Error Handler
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Malformed Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Email not verified" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Missing Firebase App Object" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or Invalid Firebase App Object")
			}

			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Token")

		}
	}

	return cfg
}

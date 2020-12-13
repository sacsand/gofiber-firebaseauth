package gofiberfirebaseauth

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for middleware
type Config struct {

	// New firebase authntication object
	// Mandatory. Default: nil
	FirebaseApp *firebase.App

	// Ignore urls array
	// Optional. Default: nil
	IgnoreUrls []string

	// Skip Email Check.
	// Optional. Default: nil
	CheckEmailVerified bool

	// Ignore email verification for these routes
	// Optional. Default: nil
	CheckEmailVerifiedIgnoredUrls []string

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Authorizer defines a function which authenticate the Authorization token and return the authenticated token
	// Optional. Default: nil
	Authorizer func(string, string) (*auth.Token, error)

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: nil
	ErrorHandler fiber.ErrorHandler

	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey string
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
	ContextKey:                    "",
}

// Initializer
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.ContextKey == "" {
		cfg.ContextKey = "user"
	}

	// Check Mandatory FirebaseApp is provided
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
		cfg.Authorizer = func(IDToken string, CurrentURL string) (*auth.Token, error) {
			if cfg.FirebaseApp == nil {
				return nil, errors.New("Missing Firebase App Object")
			}
			client, err := cfg.FirebaseApp.Auth(context.Background())
			// Verify IDToken
			token, err := client.VerifyIDToken(context.Background(), IDToken)

			// Throw error for bad token
			if err != nil {
				return nil, errors.New("Malformed Token")
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

				if checkEmail {
					// Claim email_verified from token
					if !token.Claims["email_verified"].(bool) {
						return nil, errors.New("Email not verified")
					}
				}
			}

			return token, nil
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

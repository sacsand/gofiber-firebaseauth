package gofiberfirebaseauth

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for goFiber middleware
type Config struct {
	Next func(c *fiber.Ctx) bool

	Authorizer func(string) (bool, error)

	Unauthorized fiber.Handler

	// Skip Email Check.
	// Optional. Default: nil
	CheckEmailVerified bool

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	FirebaseApp *firebase.App

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccesHandler fiber.Handler

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
	Next:               nil,
	IgnoreUrls:         nil,
	Authorizer:         nil,
	Unauthorized:       nil,
	CheckEmailVerified: false,
}

// Initialize the gofiber
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	if cfg.Authorizer == nil {
		cfg.Authorizer = func(IDToken string) (bool, error) {
			client, err := cfg.FirebaseApp.Auth(context.Background())
			// verify idTo
			token, err := client.VerifyIDToken(context.Background(), IDToken)
			// fmt.Println(err)
			if cfg.CheckEmailVerified {
				// check email is verified
				if token.Claims["email_verified"].(bool) {
					return false, errors.New("Email not verified")
				}
			}

			if err != nil {
				return false, errors.New("Malformed Token")
			}

			return true, nil
		}
	}

	// if cfg.Unauthorized == nil {
	// 	cfg.Unauthorized = func(c *fiber.Ctx) error {
	// 		//c.Set(fiber.HeaderWWWAuthenticate, "basic realm="+cfg.Realm)
	// 		return c.SendStatus(fiber.StatusUnauthorized)
	// 	}
	// }

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Malformed Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Token")

		}
	}

	// Note :: Removed desualt configuration since the go envirment complication
	// if cfg.FirebaseApp == nil {
	// 	// If the user has passed an initialized firebase app, use that
	// 	// or initialize one using the serviceAccount object.
	// 	opt := option.WithCredentialsFile("fireauth-firebase-adminsdk.json")
	// 	app, err := firebase.NewApp(context.Background(), nil, opt)

	// 	if err != nil {
	// 		log.Fatalf("error getting Auth client: %v\n", err)
	// 	}
	// 	cfg.FirebaseApp = app
	// }

	return cfg
}

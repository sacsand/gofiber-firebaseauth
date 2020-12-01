package gofiberfirebaseauth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// Config defines gofirebaseauth
type Config struct {
	Next func(c *fiber.Ctx) bool

	Authorizer func(string) bool

	Unauthorized fiber.Handler

	CheckEmailVerified bool

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	FirebaseApp *firebase.App

	SuccesHandler fiber.Handler

	ErrorHandler fiber.ErrorHandler

	IgnoreUrls []string

	CheckEmailVerifiedIgnoredUrls []string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:               nil,
	IgnoreUrls:         nil,
	Authorizer:         nil,
	Unauthorized:       nil,
	CheckEmailVerified: false,
	FirebaseApp:        nil,
}

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
		cfg.Authorizer = func(IDToken string) bool {
			client, err := cfg.FirebaseApp.Auth(context.Background())
			// verify idTo
			token, err := client.VerifyIDToken(context.Background(), IDToken)

			if cfg.CheckEmailVerified {
				log.Print(token.Claims["email_verified"].(bool))
				// check email is verified
				if token.Claims["email_verified"].(bool) {
					return false
				}
			}

			if err != nil {
				// og.Fatalf("error verifying ID token: %v\n", err)
				return false
			}

			return false
		}
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			//c.Set(fiber.HeaderWWWAuthenticate, "basic realm="+cfg.Realm)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	if cfg.FirebaseApp == nil {
		// If the user has passed an initialized firebase app, use that
		// or initialize one using the serviceAccount object.
		opt := option.WithCredentialsFile("fireauth-firebase-adminsdk.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)

		if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
		}
		cfg.FirebaseApp = app
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			// c.Set(fiber.HeaderWWWAuthenticate, "basic realm="+cfg.Realm)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

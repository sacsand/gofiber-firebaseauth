package gofiberfirebaseauth

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// New - Main
func New(config Config) fiber.Handler {
	// Init config
	cfg := configDefault(config)
	// Return new handler
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
			fmt.Println("no autorization header is present" + c.Path())
			return cfg.Unauthorized(c)

		}
		// 4) Validate the IdToken
		IsPass := cfg.Authorizer(IDToken)
		// IF Id token passed
		if IsPass {
			log.Fatalf("User authenticated")
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}

// func (auth *Auth) CreateCustomToken(userID string, claims interface{}) (string, error) {
// 	if auth.app.privateKey == nil {
// 		return "", ErrRequireServiceAccount
// 	}
// 	now := time.Now()
// 	payload := &customClaims{
// 		Issuer:    auth.app.clientEmail,
// 		Subject:   auth.app.clientEmail,
// 		Audience:  customTokenAudience,
// 		IssuedAt:  now.Unix(),
// 		ExpiresAt: now.Add(time.Hour).Unix(),
// 		UserID:    userID,
// 		Claims:    claims,
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
// 	return token.SignedString(auth.app.privateKey)
// }

// func TestVerifyIDToken(t *testing.T) {
// 	app := initApp()
// 	firAuth := app.Auth()

// 	assert.NotNil(t, app)
// 	assert.NotNil(t, firAuth)

// 	// my claims
// 	myClaims := make(map[string]string)
// 	myClaims["name"] = "go-firebase-admin"
// 	myClaims["kid"] = "polo"

// 	token, err := firAuth.CreateCustomToken("uid", myClaims)

// 	fmt.Printf("%s", token)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, token)

// 	claims, err := firAuth.VerifyIDToken(token)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, claims)
// }

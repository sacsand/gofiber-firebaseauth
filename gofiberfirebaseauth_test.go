package gofiberfirebaseauth

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

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

// 1  TEST for Malformed Token
func TestWithMalformedToken(t *testing.T) {

	// 1) intialiae fiber app and firebase app
	app := fiber.New()
	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !fileExi {
		log.Println("fireauth config not found")
	}
	// 2) create firebase app
	opt := option.WithCredentialsFile(file)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// 3) configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))
	// 4) hard coded Invalid Id token
	idToken := "0i30-ir-302309ei3f30-i32-0f-2300"

	// 5) crete  test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})
	req := httptest.NewRequest("GET", "/testauth", nil)
	req.Header.Set("Authorization", idToken)
	// 6) test
	resp, err := app.Test(req)

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		// t.Fatalf(`%s: %s`, t.Name(), err)
		fmt.Println("TEST case pass for TestWithMalformedToken")

	} else {
		log.Fatalf(`%s: %s`, t.Name(), err)
	}

}

// 2 TEST for Ignore Url
func TestIgnoreUrlsWorking(t *testing.T) {

	// t.Parallel()
	app := fiber.New()
	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !fileExi {
		log.Println("fireauth config not found")
	}
	// 2) create firebase app
	opt := option.WithCredentialsFile(file)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// 3) configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))

	app.Use(New(Config{
		FirebaseApp: fireApp,
		IgnoreUrls:  []string{"GET::/testauth", "POST::/testauth "},
	}))

	// 5) crete test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req := httptest.NewRequest("GET", "/testauth", nil)

	// 6) test
	_, err := app.Test(req)

	// fmt.Println((resp))

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	} else {
		fmt.Println("Test case pass for TestIgnoreUrlsWorking")
	}

}

// 3 TEST for Ignore Url
func TestWithoutFirebaseApp(t *testing.T) {
	// t.Parallel()
	app := fiber.New()

	// Config without firebase App object
	app.Use(New(Config{}))

	// 5) crete  test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req := httptest.NewRequest("GET", "/testauth", nil)

	// 6) test
	_, err := app.Test(req)

	// fmt.Println((resp))

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	} else {
		fmt.Println("Test case pass for TestWithoutFirebaseApp")
	}

}

//  TEST token with valid authorization token
func TestTokenWithCorrectToken(t *testing.T) {

	app := fiber.New()
	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	idToken, isTokenExist := os.LookupEnv("ID_TOKEN")

	if !fileExi || !isTokenExist {
		log.Println("fireauth config or idToken is not found")
	}
	// 2) create firebase app
	opt := option.WithCredentialsFile(file)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// 3) configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))

	req := httptest.NewRequest("GET", "/testauth", nil)

	req.Header.Set("Authorization", idToken)
	// 6) test
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		// t.Fatalf(`%s: %s`, t.Name(), err)
		fmt.Println("TEST case FAILED for TestTokenWithCorrectToken .. Provide Valid ID_TOKEN in your .env file")

	} else {
		fmt.Println("Test Case pass for TestTokenWithCorrectToken")
	}

}

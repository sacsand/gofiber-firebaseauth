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

// TEST for Malformed Token
func TestWithMalformedToken(t *testing.T) {

	t.Parallel()
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

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		// t.Fatalf(`%s: %s`, t.Name(), err)
		fmt.Println("TEST case pass for Malformed Token Check")

	} else {
		log.Fatalf(`%s: %s`, t.Name(), err)
	}

}

// TEST for Ignore Url
func TestIgnoreUrlsWorking(t *testing.T) {

	t.Parallel()
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
		fmt.Println("TEST case pass for IgnoreUrl check")
	}

}

// TEST for Ignore Url
func TestWithoutFirebaseApp(t *testing.T) {

	t.Parallel()
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
		fmt.Println("TEST case pass for No FirebaseApp")
	}

}

// // Test api passing a firebase App

// func TestIgnoreUrls(t *testing.T) {

// 	t.Parallel()

// 	app := fiber.New()
// 	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
// 	if !fileExi {
// 		log.Println("fireauth config not found")
// 	}
// 	// 2) create firebase app
// 	opt := option.WithCredentialsFile(file)
// 	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

// 	// 3) configure the gofiberfirebaseauth
// 	app.Use(New(Config{
// 		FirebaseApp: fireApp,
// 		IgnoreUrls:  []string{"GET /testauth", "POST /testauth "},
// 	}))
// 	// 4) get id token from env
// 	idToken, idExist := os.LookupEnv("ID_TOKEN")
// 	if !idExist {
// 		log.Println("no id token found")
// 	}

// 	// 5) Create  tests route
// 	app.Get("/testauth", func(c *fiber.Ctx) error {
// 		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
// 		return c.SendString(msg) // => Hello john 👋!
// 	})

// 	app.Post("/testauth", func(c *fiber.Ctx) error {
// 		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
// 		return c.SendString(msg) // => Hello john 👋!
// 	})

// 	req1 := httptest.NewRequest("GET", "/testauth", nil)
// 	req2 := httptest.NewRequest("GET", "/testauth", nil)
// 	req1.Header.Set("Authorization", idToken)
// 	req2.Header.Set("Authorization", idToken)
// 	// 6) test
// 	app.Test(req1)
// 	app.Test(req2)

// 	utils.AssertEqual(t, nil)

// }

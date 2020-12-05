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
	"gorm.io/gorm/utils"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// func TestMiddleware(t *testing.T) {

// 	// t.Parallel()
// 	app := fiber.New()

// 	// slice of array
// 	app.Use(New(Config{
// 		IgnoreUrls: []string{
// 			"GET /test/:id", "POST /test"},
// 	}))

// 	idToken, idExist := os.LookupEnv("ID_TOKEN")

// 	if !idExist {
// 		log.Printf("no id token found")
// 	}

// 	req := httptest.NewRequest("GET", "/testauth/:id", nil)

// 	req.Header.Add("Authorization", idToken)

// 	_, err := app.Test(req)

// 	utils.AssertEqual(t, nil, err)

// 	// t.FailNow()

// 	// body, err := ioutil.ReadAll(resp.Body)

// 	//type

// }

// // Test api passing a firebase App
func TestPassFireBaseObject(t *testing.T) {

	// t.Parallel()
	app := fiber.New()
	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !fileExi {
		log.Printf("fireauth config not found")
	}
	// 2) create firebase app
	opt := option.WithCredentialsFile(file)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// 3) configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))
	// 4) get id token from env
	idToken, idExist := os.LookupEnv("ID_TOKEN")
	if !idExist {
		log.Printf("no id token found")
	}

	// 5) crete  test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})
	req := httptest.NewRequest("GET", "/testauth", nil)
	req.Header.Set("Authorization", idToken)
	// 6) test
	app.Test(req)

	utils.AssertEqual(t, nil)

}

func TestIgnoreUrls(t *testing.T) {

	t.Parallel()

	app := fiber.New()
	file, fileExi := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !fileExi {
		log.Printf("fireauth config not found")
	}
	// 2) create firebase app
	opt := option.WithCredentialsFile(file)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// 3) configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
		IgnoreUrls:  []string{"GET /testauth", "POST /testauth "},
	}))
	// 4) get id token from env
	idToken, idExist := os.LookupEnv("ID_TOKEN")
	if !idExist {
		log.Printf("no id token found")
	}

	// 5) Create  tests route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	app.Post("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req1 := httptest.NewRequest("GET", "/testauth", nil)
	req2 := httptest.NewRequest("GET", "/testauth", nil)
	req1.Header.Set("Authorization", idToken)
	req2.Header.Set("Authorization", idToken)
	// 6) test
	app.Test(req1)
	app.Test(req2)

	utils.AssertEqual(t, nil)

}

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

	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))

	idToken, idExist := os.LookupEnv("ID_TOKEN")

	if !idExist {
		log.Printf("no id token found")
	}

	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req := httptest.NewRequest("GET", "/testauth", nil)

	req.Header.Set("Authorization", idToken)

	app.Test(req)

	utils.AssertEqual(t, nil)

}

// func Test_RoutesIgnore(t *testing.T) {

// 	// t.Parallel()
// 	app := fiber.New()

// 	// slice of array
// 	app.Use(New(Config{
// 		IgnoreUrls: []string{
// 			"GET /firetest", "POST /test", "GET /auth/test"},
// 	}))
// 	// TODO add support to ignore url with params
// 	// Still can ignore by putting the parms with url before the url

// 	app.Get("/auth/:x", func(c *fiber.Ctx) error {
// 		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
// 		return c.SendString(msg) // => Hello john 👋!
// 	})

// 	h := app.Handler()

// 	idToken, idExist := os.LookupEnv("ID_TOKEN")

// 	if !idExist {
// 		log.Printf("no id token found")
// 	}

// 	fctx := &fasthttp.RequestCtx{}
// 	fctx.Request.Header.SetMethod("GET")
// 	fctx.Request.SetRequestURI("/auth/test")
// 	fctx.Request.Header.Set(fiber.HeaderAuthorization, idToken) // john:doe

// 	h(fctx)

// 	// m := make(map[int]string)

// 	// req := httptest.NewRequest("GET", "/sachin", nil)

// 	// req.Header.Add("Authorization", idToken)

// 	// _, err := app.Test(req)

// 	// utils.AssertEqual(t, nil, err)

// 	// t.FailNow()

// 	// body, err := ioutil.ReadAll(resp.Body)

// 	//type

// }

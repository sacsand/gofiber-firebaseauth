package gofiberfirebaseauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Global varible for IDToken
var IDToken string

func init() {

	// loads values from .env into the system
	localDev := os.Getenv("STAGE") == ""

	if localDev {
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		}
	}
	// Get idToken form firebase and save globally
	getIDToken()
}

// Get Id token using Firebase Auth Rest API https://firebase.google.com/docs/reference/rest/auth
func getIDToken() {

	// curl 'https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=[API_KEY]' \
	// -H 'Content-Type: application/json' \
	// --data-binary '{"email":"[user@example.com]","password":"[PASSWORD]","returnSecureToken":true}'

	testUserEmail, emailExit := os.LookupEnv("TEST_USER_EMAIL")
	testUserPassword, passExit := os.LookupEnv("TEST_USER_PASSWORD")

	if !emailExit || !passExit {
		log.Println("Please provide TEST_USER_EMAIL and TEST_USER_PASSWORD")
	}

	type Payload struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ReturnSecureToken bool   `json:"returnSecureToken"`
	}

	data := Payload{
		Email:             testUserEmail,
		Password:          testUserPassword,
		ReturnSecureToken: true,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	// body := strings.NewReader(`{"email":"skodagoda@apm.mc","password":"sac1234","returnSecureToken":true}`)
	body := bytes.NewReader(payloadBytes)

	webAPIKey, keyExitx := os.LookupEnv("WEB_API_KEY")

	if !keyExitx {
		log.Println("WEB_API_KEY is not configured.Please add the WEB_API_KEY to your .env")
	}

	req, err := http.NewRequest("POST", "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+webAPIKey+"", body)
	if err != nil {
		log.Println("Error generating idToken")
	}

	req.Header.Set("content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error generating IDToken")
	}

	// type response interface {
	// }

	type response struct {
		IDToken string `json:"idToken"`
	}

	defer resp.Body.Close()
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Getting IDToken")
	}
	// fmt.Print(bodyResponse)
	var Response response
	json.Unmarshal(bodyResponse, &Response)
	//fmt.Println(Response.IDToken)
	IDToken = Response.IDToken
}

// 1  TEST for Malformed Token
func TestWithMalformedToken(t *testing.T) {

	// intialiae fiber app and firebase app
	app := fiber.New()
	serviceAccountJSON := os.Getenv("SERVICE_ACCOUNT_JSON")
	// if !fileExi {
	// 	log.Println("fireauth config not found")
	// }
	// create firebase app
	opt := option.WithCredentialsFile(serviceAccountJSON)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))
	// hard coded Invalid Id token
	idToken := "0i30-ir-302309ei3f30-i32-0f-2300"

	// crete  test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})
	req := httptest.NewRequest("GET", "/testauth", nil)
	req.Header.Set("Authorization", idToken)
	// test
	resp, err := app.Test(req)

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		fmt.Println("TEST case pass for TestWithMalformedToken")
	} else {
		log.Fatalf(`%s: %s`, t.Name(), err)
	}

}

// 2 TEST for Ignore Url
func TestIgnoreUrlsWorking(t *testing.T) {

	// t.Parallel()
	app := fiber.New()
	serviceAccountJSON, fileExi := os.LookupEnv("SERVICE_ACCOUNT_JSON")
	if !fileExi {
		log.Println("fireauth config not found")
	}
	// create firebase app
	opt := option.WithCredentialsFile(serviceAccountJSON)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))

	app.Use(New(Config{
		FirebaseApp: fireApp,
		IgnoreUrls:  []string{"GET::/testauth", "POST::/testauth "},
	}))

	// crete test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req := httptest.NewRequest("GET", "/testauth", nil)

	// test
	_, err := app.Test(req)

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	} else {
		fmt.Println("Test case pass for TestIgnoreUrlsWorking")
	}

}

// 3 TEST for FirebaseApp
func TestWithoutFirebaseApp(t *testing.T) {
	// t.Parallel()
	app := fiber.New()

	// Config without firebase App object
	app.Use(New(Config{}))

	// crete  test route
	app.Get("/testauth", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	req := httptest.NewRequest("GET", "/testauth", nil)

	// test
	_, err := app.Test(req)

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	} else {
		fmt.Println("Test case pass for TestWithoutFirebaseApp")
	}

}

// 4 TEST token with valid authorization token
func TestTokenWithCorrectToken(t *testing.T) {

	app := fiber.New()

	serviceAccountJSON, fileExi := os.LookupEnv("SERVICE_ACCOUNT_JSON")
	if !fileExi {
		log.Println("fireauth config not found")
	}

	// create firebase app
	opt := option.WithCredentialsFile(serviceAccountJSON)
	fireApp, _ := firebase.NewApp(context.Background(), nil, opt)

	// configure the gofiberfirebaseauth
	app.Use(New(Config{
		FirebaseApp: fireApp,
	}))

	req := httptest.NewRequest("GET", "/testauth", nil)

	req.Header.Set("Authorization", IDToken)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		fmt.Println("TEST case FAILED for TestTokenWithCorrectToken ")
	} else {
		fmt.Println("Test Case pass for TestTokenWithCorrectToken")
	}

}

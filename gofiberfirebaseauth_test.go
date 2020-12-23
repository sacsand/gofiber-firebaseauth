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

	localDev := os.Getenv("STAGE") == ""
	// loads values from .env into the system
	if localDev { // no in build, only in local
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		}
	}
	// Get idToken form firebase and save globally
	getIDToken()
}

/**
 *
 *	Helper Functions
 *
 */

// Get Idtoken by calling Firebase Auth Rest API. DOCS:: https://firebase.google.com/docs/reference/rest/auth
func getIDToken() {

	// curl 'https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=[API_KEY]' \
	// -H 'Content-Type: application/json' \
	// --data-binary '{"email":"[user@example.com]","password":"[PASSWORD]","returnSecureToken":true}'

	// Load envirment variable
	testUserEmail, emailExit := os.LookupEnv("TEST_USER_EMAIL")
	testUserPassword, passExit := os.LookupEnv("TEST_USER_PASSWORD")
	if !emailExit || !passExit {
		log.Println("Please provide TEST_USER_EMAIL and TEST_USER_PASSWORD")
	}

	webAPIKey, keyExitx := os.LookupEnv("WEB_API_KEY")
	if !keyExitx {
		log.Println("WEB_API_KEY is not configured.Please add the WEB_API_KEY to your .env")
	}

	// preparing payload
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
		log.Println("Error getting idToken")
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+webAPIKey+"", body)
	if err != nil {
		log.Println("Error generating idToken")
	}

	req.Header.Set("content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error generating IDToken")
	}

	type response struct {
		IDToken string `json:"idToken"`
	}

	defer resp.Body.Close()
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Getting IDToken")
	}

	var Response response
	json.Unmarshal(bodyResponse, &Response)
	IDToken = Response.IDToken
}

// Create FirebaseAuth
func CreateFirebaseAuthApp() (*firebase.App, error) {

	serviceAccountJSON, fileExi := os.LookupEnv("SERVICE_ACCOUNT_JSON")

	if !fileExi {
		log.Println("fireauth config not found")
	}

	// Create a firebase app
	opt := option.WithCredentialsFile(serviceAccountJSON)
	fireApp, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		return nil, err
	}

	return fireApp, nil

}

/**
*
*	TEST CASES
*
 */

// 1  TEST for Malformed Token
func TestWithMalformedToken(t *testing.T) {

	// intialiae fiber app and firebase app
	app := fiber.New()

	fireApp, err := CreateFirebaseAuthApp()

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

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

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

	if resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnauthorized {
		// fmt.Println("TEST case pass for TestWithMalformedToken")
	} else {
		log.Fatalf(`%s: %s`, t.Name(), err)
	}

}

// 2 TEST for Ignore Url
func TestIgnoreUrlsWorking(t *testing.T) {

	// t.Parallel()
	app := fiber.New()

	// Create firebase app
	fireApp, errf := CreateFirebaseAuthApp()

	if errf != nil {
		t.Fatalf(`%s: %s`, t.Name(), errf)
	}

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
	}

}

// 4 TEST token with valid authorization token
func TestTokenWithCorrectToken(t *testing.T) {

	app := fiber.New()

	fireApp, err := CreateFirebaseAuthApp()

	if err != nil {
		t.Fatalf(`%s: %s`, t.Name(), err)
	}

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
	}

}

# Go Fiber Firebase Auth Middleware 
[![CircleCI](https://circleci.com/gh/sacsand/gofiber-firebaseauth.svg?style=shield)](https://circleci.com/gh/sacsand/gofiber-firebaseauth)

Authenticate your endpoints with [Firebase Authentication ](https://github.com/LeafyCode/express-firebase-auth/).

gofiber-firebaseauth is inspired by npm package [express-firebase-auth](https://github.com/LeafyCode/express-firebase-auth/) .
#### Note1
 This package is designed to work with [Go Fiber Framework](https://github.com/gofiber/fiber) (Express inspired go framework)
# Features
  - Authenticate the user using Firebase before running the function.
  - Ability to skip authentication on public API endpoints.
  
# Installing / Getting started

In your gofiber app
```sh
$ go get -u github.com/gofiber/fiber/v2
$ go get github.com/sacsand/gofiber-firebaseauth
```

# Configure 
In your app import the middleware

```go
import (
  "github.com/gofiber/fiber/v2"
  "github.com/sacsand/gofiber-firebaseauth"
)
```
Provide a minimal config
```go
    // Provide a minimal config
    app.Use(gofiberfirebaseauth.New(Config{
    	FirebaseApp: FirebaseApp,
    }))
```
Or extend your config for customization
```go

    // Or extend your config for customization
	app.Use(gofiberfirebaseauth.New(Config{
	
	// New firebase authentication object
	// Mandatory. Default: nil
	FirebaseApp:  FirebaseApp

	// Ignore urls array - Format = "{METHOD} follwed by :: then /{route}"
	// Optional. Default: nil
	IgnoreUrls : []string{"GET::/login","POST::/create-user"}

	// Skip Email Check.
	// Optional. Default: nil
	CheckEmailVerified : true

	// Ignore email verification for these routes
	// Optional. Default: nil
	CheckEmailVerifiedIgnoredUrls :  []string{"GET::/login","POST::/create-user"}

	// Authorizer defines a function which authenticates the Authorization token and returns 
	// the authenticated token
	// Optional. Default: nil
	Authorizer: func(IDToken string, CurrentURL string) (*auth.Token, error){
        // create your own authentication here 
	// this returns the firebase id token
		return token, nil
	},
	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey : "authUser"
    }))
```
Use user in your fiber app
```go
func Handler(ctx *fiber.Ctx) error {
	// Get user stored in context
	// Default: user
	currentUser := ctx.Locals("user").(gofiberfirebaseauth.User)
	fmt.Println(currentUser)
	fmt.Println(currentUser.Email)

}
```

All available configuration

| Option                          |                                                                                                                                                                                            Value | Config type                                          |
| ------------------------------- | -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------- |
| `FirebaseApp`                   |                                                                       ([**Note1**](#note2)) An initialized firebase app. [Refer to Firebase setup](https://firebase.google.com/docs/admin/setup) | FirebaseApp *firebase.App                            |
| `IgnoredUrls`                   |                                                                                                    ([**Note1**](#note3))(*Optional*) An array of URLs where you need to skip the authentication. | IgnoreUrls []string                                  |
| `CheckEmailVerified`            |                                                                                   (*Optional*) (Default: **false**) If set to **true**, only users with a verified email will be allowed access. | CheckEmailVerified bool                              |
| `CheckEmailVerifiedIgnoredUrls` |                                                                                                                   (*Optional*) An array of URLs where you need to skip the email verified check. | CheckEmailVerifiedIgnoredUrls []string               |
| `Authorizer`                    | (*Optional*), Default: nil , Authorizer defines a function which authenticates the Authorization token and returns the authenticated token. Use this if you want to override token authorization | Authorizer func(string, string) (*auth.Token, error) |
| `ContextKey`                    |                                                                                               (*Optional*), Default: "user" , Context key to store user information from the token into context. | ContextKey string                                    |
| `SuccessHandler`                |                                                                                                (*Optional*), Default:nil, SuccessHandler defines a function which is executed for a valid token. | SuccessHandler fiber.Handler                         |
| `ErrorHandler`                  |                                                                                                (*Optional*), Default:nil, ErrorHandler defines a function which is executed for a invalid token. | ErrorHandler fiber.ErrorHandler                      |
                            


 #### Note2
You **must** provide already initialized `FirebaseApp` app.
You cannot initialize two firebase apps.

 #### Note3
  Ignore url accept array of string. URl format should follow below format\
  `{METHOD}::/{url}`\
  Example: \
  GET::/login\
  POST::/login
  
  IgnoreUrl only supports routes without params or query . (PR are welcome).To ignore urls with param or query,  declare the routes before the middleware declaration. 

## Developing , TESTING and environment setup

### Prerequisites
- Go 1.14 +
- Configured Firebase app and Google Service Account Credential (JSON containing admin credentials). Refer to [Firebase setup](https://firebase.google.com/docs/admin/setup) 
- Web API Key
- Sample user email and password from firebase. You can manually create a user from the firebase console.

You can get all the configurations from Firebase Console.
 
### Setting up Development Environment for testing

Clone the repo and set your firebase credentials in your .env file

```

SERVICE_ACCOUNT_JSON = "path to service account credential json"
WEB_API_KEY = 
TEST_USER_EMAIL = ""
TEST_USER_PASSWORD = ""
  	
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [link to tags on this repository](/releases).

## Style guide

 [Uber](https://github.com/uber-go/guide/blob/master/style.md ) style guide


## License

[MIT licensed](./LICENSE).

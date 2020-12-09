# Gofiber Firebase Auth Middleware 

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

Authenticate your endpoints with [Firebase Authentication ](https://github.com/LeafyCode/express-firebase-auth/).

gofiberfirebase auth is inspired by npm package [express-firebase-auth](https://github.com/LeafyCode/express-firebase-auth/) .
NOTE :: This package design to work with [Go Fiber Framework](https://github.com/gofiber/fiber) (Express inspired go framework)
# Features
  - Authenticate the user using Firebase before running the function.
  - Ability to skip authentication on public API endpoints.
  
# Installing / Getting started

In your gofiber app
```sh
$ go get -u github.com/gofiber/fiber/v2
$ go get gofiberfirebaseauth
```


# Configure 
In your app import the middleware

```go
import (
  "github.com/gofiber/fiber/v2"
  "github.com/sacsand/gofiberfirebaseauth"
)
```



```go
// Provide a minimal config
	app.Use(gofiberfirebaseauth.New(Config{
    	FirebaseApp: fireApp,
    }))

// Or extend your config for customization
	app.Use(gofiberfirebaseauth.New(Config{
	    // Firebase Authentication App Object
	    // Mandatory 
    	FirebaseApp: fireApp,
    	
    	// Ignore urls array
    	// Optional
    	IgnoreUrls: []string{"GET::/testauth", "POST::/testauth "},
    	
    	// Token authorizer
        // Optional
        Authorizer : func(IDToken string, CurrentURL string) (bool, error) {
          //  Override the Authorization here
		},
		
		// ErrorHandler defines a function which is executed for an invalid token
		// Optional
		ErrorHandler : func(c *fiber.Ctx, err error) error {
		 // Override ErrorHandler here
		},
		SuccessHandler : func(c *fiber.Ctx) error {
		 // Override SuccessHandler here
		},
		CheckEmailVerified : true,
	    CheckEmailVerifiedIgnoredUrls: []string{"GET::/testauth", "POST::/testauth "},
	}
		
    }))

```

 # TODO  add developer guide and contribution guide





# Gofiber Firebase Auth Middleware (-- Under Testing)

Authenticate your endpoints with [Firebase Authentication ](https://github.com/LeafyCode/express-firebase-auth/).

gofiberfirebase auth is inspired by npm package [express-firebase-auth](https://github.com/LeafyCode/express-firebase-auth/) .
#### Note1
 This package design to work with only [Go Fiber Framework](https://github.com/gofiber/fiber) (Express inspired go framework)
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
    	FirebaseApp: FirebaseApp,
    }))

    // Or extend your config for customization
	app.Use(gofiberfirebaseauth.New(Config{
	
	// New firebase authntication object
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

	// Authorizer defines a function which authenticate the Authorization token and return the authenticated token
	// Optional. Default: nil
	Authorizer: func(IDToken string, CurrentURL string) (*auth.Token, error){
    // create your own authentication here 
	// this retuen firbase id token
		return token, nil
	},
	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey : "authUser"
    }))

```

All available configuration

| Option                  | Value           |Config type
| -------------           |-----------------------------------------------------------------------------------------------------------:|:------
| `FirebaseApp`              | ([**Note1**](#note1)) An initialized firebase app. [Refer Firebase setup](https://firebase.google.com/docs/admin/setup)   | FirebaseApp *firebase.App
| `IgnoredUrls`           | (*Optional*) An array of URLs where you need to skip the authentication.                                                  | IgnoreUrls []string
| `CheckEmailVerified`    | (*Optional*) (Default: **false**) If set to **true**, only users with a verified email will be allowed access.            | CheckEmailVerified bool
| `CheckEmailVerifiedIgnoredUrls`    | (*Optional*) An array of URLs where you need to skip the email verified check.                                 | CheckEmailVerifiedIgnoredUrls []string
| `Authorizer`    | (*Optional*), Default: nil , Authorizer defines a function which authenticate the Authorization token and return the authenticated token. Use this if you want to oevride token authorization |Authorizer func(string, string) (*auth.Token, error)
| `ContextKey`    | (*Optional*), Default: "user" , Context key to store user information from the token into context.| ContextKey string
| `SuccessHandler`    |	(*Optional*), Default:nil, SuccessHandler defines a function which is executed for a valid token.  | SuccessHandler fiber.Handler
| `ErrorHandler`    |	(*Optional*), Default:nil, ErrorHandler defines a function which is executed for a invalid token. | ErrorHandler fiber.ErrorHandler
                            


 #### Note2
You **must** provide already initialized `FirebaseApp` app.
You cannot initialize two firebase apps.

This package adds the `user` object returned by firebase to context key `user`. You can use that inside your functions.
(you have option to configure the ContextKey. Check the configarstion available table )

user object type
```
        	type user struct {
				email, emailVerified bool
				userID               string
			}

```



## Developing and TESTING envirnment setup

### Prerequisites
-  Go 1.14 +
- Configured Firebase Authentication App   GOOGLE_APPLICATION_CREDENTIALS
- Valid ID_TOKEN , how to get https://github.com/sacsand/firebase-get-IdToken

### Setting up Dev

Clone the repo and set your firebase credentials in your .env file\  
example::\
GOOGLE_APPLICATION_CREDENTIALS = "path/firebase-admin.json"\
ID_TOKEN = "eyryursdkdjksjdks"  	


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [link to tags on this repository](/releases).

## Style guide

 [Uber](https://github.com/uber-go/guide/blob/master/style.md ) stryle guide


## License

[MIT licensed](./LICENSE).

## TODO

TODO add support for Ignore URLS to to skip route with params and quarries





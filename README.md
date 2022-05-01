# JWT Middleware For Go Fiber

This library is to make it easier for you to limit access rights, especially using Json Web Token

## Installation

To Install JWT Middaleware Run

```bash
  go get github.com/asnur/fiber_jwtware
```

## Features

- Set Cookie For Web Route
- Restricred Web Route
- Dynamic For Route Web or Api

## Usage/Examples

```golang
package main

import (
	jwtware "github.com/asnur/fiber_jwtware"
	"github.com/gofiber/fiber/v2"
)

func main() {
  app := fiber.New()

 // JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
        Type: "api" //optional
		    Secret: "YOUR SECRET KEY",
        Redirect: "/login" //optional if you use Type is web
	}))

//Restricted Route
    app.Get('/', ....)

//Other your routes
}
```

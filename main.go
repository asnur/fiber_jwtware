package fiber_jwtware

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Type         string                                     //api or web
	Next         func(c *fiber.Ctx) bool                    // Required
	Unauthorized fiber.Handler                              // middleware specfic
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error) // middleware specfic
	Secret       string                                     // middleware specfic
	Expiry       int64
	Redirect     string // optional if you use Type is Web (Not API)
}

var ConfigDefault = Config{
	Type:         "api",
	Next:         nil,
	Decode:       nil,
	Unauthorized: nil,
	Secret:       "secret",
	Expiry:       60,
	Redirect:     "/login",
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	// Set default secret if not passed
	if cfg.Secret == "" {
		cfg.Secret = ConfigDefault.Secret
	}

	// Set default expiry if not passed
	if cfg.Expiry == 0 {
		cfg.Expiry = ConfigDefault.Expiry
	}

	// this is the main jwt decode function of our middleware
	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {

			authHeader := c.Cookies("Token")

			if authHeader == "" {
				return nil, errors.New("Authorization header is required")
			}

			// we parse our jwt token and check for validity against our secret
			token, err := jwt.Parse(
				authHeader,
				func(token *jwt.Token) (interface{}, error) {
					// verifying our algo
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(cfg.Secret), nil
				},
			)

			if err != nil {
				return nil, errors.New("Error parsing token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !(ok && token.Valid) {
				return nil, errors.New("Invalid token")
			}

			if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("jwt is expired")
			}
			fmt.Println(int64(claims["exp"].(float64)), time.Now().UTC().Unix())
			fmt.Println(int64(claims["exp"].(float64)) < time.Now().UTC().Unix())
			fmt.Println(time.Unix(int64(claims["exp"].(float64)), 0), time.Now())

			return &claims, nil
		}
	}

	// Set default Unauthorized if not passed
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			if cfg.Type == "web" {
				return c.Redirect(cfg.Redirect)
			} else {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}
	}

	return cfg
}

func New(config Config) fiber.Handler {

	// For setting default config
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			fmt.Println("Midddle was skipped")
			return c.Next()
		}
		fmt.Println("Midddle was run")

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("jwtClaims", *claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}

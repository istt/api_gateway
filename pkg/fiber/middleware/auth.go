package middleware

import (
	"fmt"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// JWTSECRET hold the JWT secret for encode and decode
var JWTSECRET string

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(JWTSECRET),
		ErrorHandler: jwtError,
	})
}

// HasAuthority check if the current role has specified authorities
func HasAuthority(authorityName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if authorities, ok := claims["authorities"]; ok && (authorities != nil) {
			authorities := claims["authorities"].([]interface{})
			for _, authority := range authorities {
				if fmt.Sprint(authority) == authorityName {
					return c.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf(`Account doesn't have required authority %s to access this resource`, authorityName))
	}
}

// jwtError return error for JWT
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return fiber.NewError(fiber.StatusBadRequest, "missing or malformed authorization token")
	}
	return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired authorization token")
}

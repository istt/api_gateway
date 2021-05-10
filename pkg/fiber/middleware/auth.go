package middleware

import (
	"fmt"
	"strings"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// JWTSECRET hold the JWT secret for encode and decode
var JWTSECRET string

const AUTHORITIES_KEY = "auth"
const FIBER_CONTEXT_KEY = "user"

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
		user := c.Locals(FIBER_CONTEXT_KEY).(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if authorities, ok := claims[AUTHORITIES_KEY]; ok && (authorities != nil) {
			authorities := strings.Split(claims[AUTHORITIES_KEY].(string), ",")
			for _, authority := range authorities {
				if fmt.Sprint(authority) == authorityName {
					return c.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "missing required authority to access this resource")
	}
}

// jwtError return error for JWT
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return fiber.NewError(fiber.StatusBadRequest, "missing or malformed authorization token")
	}
	return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired authorization token")
}

// GetCurrentUserLogin return the current login name for this particular user
func GetCurrentUserLogin(c *fiber.Ctx) (string, error) {
	ctx := c.Locals(FIBER_CONTEXT_KEY)
	if ctx == nil {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "missing JWT information")
	}
	token, ok := ctx.(*jwt.Token)
	if !ok {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "invalid JWT information")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "invalid JWT information")
	}
	return fmt.Sprint(claims["sub"]), nil
}

package rest

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	var input shared.LoginVM

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	ud, err := instances.UserService.GetUserByUsername(c.Context(), input.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid login")
	}
	if !ud.Activated {
		return fiber.NewError(fiber.StatusExpectationFailed, "Account is not activated")
	}
	if !instances.UserService.CheckPasswordHash(input.Password, ud.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = ud.Login
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["authorities"] = ud.Authorities

	t, err := token.SignedString([]byte(app.Config.String("security.jwt-secret")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to generate token")
	}
	return c.JSON(fiber.Map{"id_token": t})
}

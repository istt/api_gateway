package rest

import (
	"fmt"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/pkg/fiber/middleware"
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
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("unable to find data: %s | %s", input.Username, err.Error()))
	}
	if ud == nil {
		return fiber.NewError(fiber.StatusForbidden, "invalid username or password")
	}
	if !ud.Activated {
		return fiber.NewError(fiber.StatusExpectationFailed, "account is not activated")
	}
	if !instances.UserService.CheckPasswordHash(input.Password, ud.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	// mimic jhipster claims
	claims := jwt.MapClaims{
		"sub":  ud.Login,
		"auth": strings.Join(ud.Authorities, ","),
	}

	if input.RememberMe {
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t, err := token.SignedString([]byte(middleware.JWTSECRET))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to generate token")
	}
	return c.JSON(fiber.Map{"id_token": t})
}

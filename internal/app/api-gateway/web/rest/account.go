package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/pkg/fiber/authjwt/utils"
	"github.com/istt/api_gateway/pkg/fiber/instances"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

// GetAccount implement api endpoint
func GetAccount(c *fiber.Ctx) error {
	log.Print("access account")
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := instances.UserService.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	return c.JSON(account.UserDTO)
}

// SaveAccount implement api endpoint
func SaveAccount(c *fiber.Ctx) error {
	var updatedInfo shared.UserDTO
	if err := c.BodyParser(&updatedInfo); err != nil {
		return err
	}
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := instances.UserService.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	account.UserDTO = updatedInfo
	if err := instances.UserService.SaveAccount(c.Context(), account); err != nil {
		return err
	}
	return c.JSON(account.UserDTO)
}

// ChangePassword implement api endpoint
func ChangePassword(c *fiber.Ctx) error {
	var input shared.PasswordChangeDTO
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := instances.UserService.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	if !instances.UserService.CheckPasswordHash(input.CurrentPassword, account.Password) {
		return fiber.NewError(fiber.StatusExpectationFailed, "current password does not match")
	}

	hash, err := instances.UserService.HashPassword(input.NewPassword)
	if err != nil {
		return err

	}

	account.Password = hash
	if err := instances.UserService.SaveAccount(c.Context(), account); err != nil {
		return err
	}
	return c.JSON(account.UserDTO)
}

// FinishPasswordReset implement api endpoint
func FinishPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// RequestPasswordReset implement api endpoint
func RequestPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// ActivateAccount implement api endpoint
func ActivateAccount(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// IsAuthenticated implement api endpoint
func IsAuthenticated(c *fiber.Ctx) error {
	return Login(c)
}

// RegisterAccount implement api endpoint
func RegisterAccount(c *fiber.Ctx) error {
	var user shared.UserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	exists, err := instances.UserService.GetUserByUsername(c.Context(), user.Login)
	if err != nil {
		return err
	}
	if exists != nil {
		return fiber.NewError(fiber.StatusConflict, "user exists")
	}

	hash, err := instances.UserService.HashPassword(utils.RandomString(8))
	if err != nil {
		return err
	}
	newAccount := shared.ManagedUserDTO{
		UserDTO:  user,
		Password: hash,
	}
	log.Printf("saving data to database: %+v", newAccount)
	if err := instances.UserService.SaveAccount(c.Context(), &newAccount); err != nil {
		return err
	}
	return c.JSON(user)
}

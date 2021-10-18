package rest

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/istt/api_gateway/pkg/fiber/authjwt/consts"
	"github.com/istt/api_gateway/pkg/fiber/authjwt/utils"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

type AccountResource interface {
	// + Jhipster endpoint for ROLE_USER
	GetAccount(c *fiber.Ctx) error           // getAccount
	SaveAccount(c *fiber.Ctx) error          // saveAccount
	ChangePassword(c *fiber.Ctx) error       // ChangePassword
	FinishPasswordReset(c *fiber.Ctx) error  // finishPasswordReset
	RequestPasswordReset(c *fiber.Ctx) error // requestPasswordReset
	ActivateAccount(c *fiber.Ctx) error      // activateAccount
	Login(c *fiber.Ctx) error                // isAuthenticated
	Register(c *fiber.Ctx) error             // registerAccount
}

type DefaultAccountResource struct {
	UserSvc services.UserService
}

func NewDefaultAccountResource(u services.UserService) AccountResource {
	return &DefaultAccountResource{
		UserSvc: u,
	}
}

// GetAccount implement api endpoint
func (r *DefaultAccountResource) GetAccount(c *fiber.Ctx) error {
	log.Print("access account")
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := r.UserSvc.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	return c.JSON(account.UserDTO)
}

// SaveAccount implement api endpoint
func (r *DefaultAccountResource) SaveAccount(c *fiber.Ctx) error {
	var updatedInfo shared.UserDTO
	if err := c.BodyParser(&updatedInfo); err != nil {
		return err
	}
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := r.UserSvc.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	account.UserDTO = updatedInfo
	if err := r.UserSvc.SaveAccount(c.Context(), account); err != nil {
		return err
	}
	return c.JSON(account.UserDTO)
}

// ChangePassword implement api endpoint
func (r *DefaultAccountResource) ChangePassword(c *fiber.Ctx) error {
	var input shared.PasswordChangeDTO
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	login, err := utils.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err := r.UserSvc.GetUserByUsername(c.Context(), login)
	if err != nil {
		return err
	}

	if !r.UserSvc.CheckPasswordHash(input.CurrentPassword, account.Password) {
		return fiber.NewError(fiber.StatusExpectationFailed, "current password does not match")
	}

	hash, err := r.UserSvc.HashPassword(input.NewPassword)
	if err != nil {
		return err

	}

	account.Password = hash
	if err := r.UserSvc.SaveAccount(c.Context(), account); err != nil {
		return err
	}
	return c.JSON(account.UserDTO)
}

// FinishPasswordReset implement api endpoint
func (r *DefaultAccountResource) FinishPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// RequestPasswordReset implement api endpoint
func (r *DefaultAccountResource) RequestPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// ActivateAccount implement api endpoint
func (r *DefaultAccountResource) ActivateAccount(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// Login implement api endpoint
func (r *DefaultAccountResource) Login(c *fiber.Ctx) error {
	var input shared.LoginVM

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	ud, err := r.UserSvc.GetUserByUsername(c.Context(), input.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid login")
	}
	if !ud.Activated {
		return fiber.NewError(fiber.StatusExpectationFailed, "account is not activated")
	}
	if !r.UserSvc.CheckPasswordHash(input.Password, ud.Password) {
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
	t, err := token.SignedString([]byte(consts.JWTSECRET))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to generate token")
	}
	return c.JSON(fiber.Map{"id_token": t})
}

// Register implement api endpoint
func (r *DefaultAccountResource) Register(c *fiber.Ctx) error {
	var user shared.UserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	exists, err := r.UserSvc.GetUserByUsername(c.Context(), user.Login)
	if err == nil && exists != nil {
		return fiber.NewError(fiber.StatusConflict, "user exists")
	}

	hash, err := r.UserSvc.HashPassword(utils.RandomString(8))
	if err != nil {
		return err
	}
	newUser := shared.ManagedUserDTO{
		UserDTO:     user,
		Password:    hash,
		CreatedBy:   user.Login,
		CreatedDate: time.Now().Format(time.RFC3339),
	}
	newUser.Activated = false
	newUser.Authorities = []string{"ROLE_USER"}
	if err := r.UserSvc.SaveAccount(c.Context(), &newUser); err != nil {
		return err
	}
	return c.JSON(user)
}

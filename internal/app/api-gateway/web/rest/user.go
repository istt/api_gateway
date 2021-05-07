package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

// GetAllUser get a user
func GetAllUser(c *fiber.Ctx) error {
	rows := make([]shared.UserDTO, 0)
	c.Set("X-Total-Count", "0")
	return c.JSON(rows)
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id parameter")
	}
	var user shared.UserDTO
	return c.JSON(user)
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	// if err := c.BodyParser(&user); err != nil {
	// 	return err
	// }
	// resp, _ := app.EtcdClient.Get(c.Context(), "U-"+user.Login)
	// if resp.Count > 0 {
	// 	return fiber.ErrConflict
	// }
	// if strings.TrimSpace(user.Password) == "" {
	// 	user.Password = RandStringBytesMaskImprSrcUnsafe(8)
	// 	c.Response().Header.Add("X-Password", user.Password)
	// }
	// hash, err := services.HashPassword(user.Password)
	// if err != nil {
	// 	return err
	// }

	// user.Password = hash
	// jsondata, er := json.Marshal(user)
	// if er != nil {
	// 	return er
	// }
	// if _, err := app.EtcdClient.Put(c.Context(), fmt.Sprintf("U-%s", user.Login), string(jsondata)); err != nil {
	// 	return err
	// }
	return c.JSON(user)

}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	// if err := c.BodyParser(&user); err != nil {
	// 	return err
	// }
	// resp, _ := app.EtcdClient.Get(c.Context(), "U-"+user.Login)
	// if resp.Count != 1 {
	// 	return fiber.NewError(fiber.StatusNotFound, "Unable to find user with name "+user.Login)
	// }
	// if err := json.Unmarshal(resp.Kvs[0].Value, &existsUser); err != nil {
	// 	return err
	// }
	// if strings.TrimSpace(user.Password) == "" {
	// 	user.Password = existsUser.Password
	// } else {
	// 	hash, err := services.HashPassword(user.Password)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	user.Password = hash
	// }

	// jsondata, er := json.Marshal(user)
	// if er != nil {
	// 	return er
	// }
	// if _, err := app.EtcdClient.Put(c.Context(), fmt.Sprintf("U-%s", user.Login), string(jsondata)); err != nil {
	// 	return err
	// }
	return c.JSON(user)

}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}

// GetAuthorities return list of authorities
func GetAuthorities(c *fiber.Ctx) error {
	return c.JSON(app.Config.Strings("security.authorities"))
}

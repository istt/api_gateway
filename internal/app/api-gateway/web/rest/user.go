package rest

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/pkg/common/utils"
	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
	"github.com/istt/api_gateway/pkg/fiber/middleware/filter/mgo"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUser get a user
func GetAllUser(c *fiber.Ctx) error {
	predicate, ok := c.Locals("filter").(filter.Filter)
	if ok {
		return GetAllUserWithPredicate(c, predicate)
	}
	cur, err := app.MongoDB.Collection("user").Find(c.Context(), bson.D{})
	if err != nil {
		return err
	}
	entities := make([]shared.UserDTO, 0)
	if err := cur.All(c.Context(), &entities); err != nil {
		return err
	}
	c.Set("X-Total-Count", "0")
	return c.JSON(entities)
}

func GetAllUserWithPredicate(c *fiber.Ctx, predicate filter.Filter) error {
	filterMongo := mgo.NewFilterMongo(&predicate)
	filterMap, findOptions, err := filterMongo.MarshalBSON()
	if err != nil {
		return err
	}
	cur, err := app.MongoDB.Collection("user").Find(c.Context(), filterMap, findOptions)
	if err != nil {
		return err
	}
	entities := make([]shared.UserDTO, 0)
	if err := cur.All(c.Context(), &entities); err != nil {
		return err
	}
	return c.JSON(entities)
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id parameter")
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var entity shared.ManagedUserDTO
	if err := app.MongoDB.Collection("user").FindOne(c.Context(), bson.M{"_id": objID}).Decode(&entity); err != nil {
		return err
	}
	return c.JSON(entity.UserDTO)
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	if user.Id != "" {
		return fiber.NewError(fiber.StatusBadRequest, "new user cannot have ID")
	}

	if strings.TrimSpace(user.Password) == "" {
		user.Password = utils.RandomString(8)
		c.Response().Header.Add("X-Password", user.Password)
	}
	hash, err := instances.UserService.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash
	if err := instances.UserService.SaveAccount(c.Context(), &user); err != nil {
		return err
	}
	return c.JSON(user.UserDTO)

}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	var user shared.ManagedUserDTO
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	existsUser, err := instances.UserService.GetUserByUsername(c.Context(), user.Login)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "username or email not found")
	}
	if strings.TrimSpace(user.Password) == "" {
		user.Password = existsUser.Password
	} else {
		hash, err := instances.UserService.HashPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = hash
	}
	if err := instances.UserService.SaveAccount(c.Context(), &user); err != nil {
		return err
	}
	return c.JSON(user)

}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.ErrNotFound
	}
	res, err := app.MongoDB.Collection("user").DeleteOne(c.Context(), bson.M{"login": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fiber.ErrNotFound
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetAuthorities return list of authorities
func GetAuthorities(c *fiber.Ctx) error {
	return c.JSON(app.Config.Strings("security.authorities"))
}

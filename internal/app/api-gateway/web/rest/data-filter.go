package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
)

func DataFilter(c *fiber.Ctx) error {
	if queryFilters, ok := c.Locals("filter").(filter.Filter); ok {
		log.Printf("%+v", queryFilters)
		return c.JSON(queryFilters)
	}
	return fiber.ErrBadRequest
}

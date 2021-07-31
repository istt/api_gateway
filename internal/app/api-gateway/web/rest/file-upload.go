package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// FileUpload let user upload a `.cdb` file then convert it into key value pair in raw text
func FileUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	if err = c.SaveFile(file, fmt.Sprintf("/tmp/%s", file.Filename)); err != nil {
		return err
	}
	return c.JSON(file.Filename)
}

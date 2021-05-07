package utils

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PagingParam Helper func to parse paging param given by query param
func PagingParam(c *fiber.Ctx) (int64, int64, int64, int64) {
	page, err := strconv.ParseInt(c.Query("page", "0"), 0, 0)
	if err != nil {
		page = 1
	}
	if page < 0 {
		page = 1
	}
	size, err := strconv.ParseInt(c.Query("size", "20"), 0, 0)
	if err != nil {
		size = 20
	}
	if size < 1 {
		size = 20
	}
	limit := (page + 1) * size
	skip := page * size
	log.Printf("page %d size %d limit %d offset %d", page, size, limit, skip)
	return page, size, limit, skip
}

package filter

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		ctxFilter := Filter{}
		dataFilters := make([]*KeyvalueOperator, 0)
		sort := Sort{}
		page := PageRequest{Page: 1, Size: 20}
		paramsMap := make(map[string][]string)
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			val, ok := paramsMap[string(key)]
			if ok {
				paramsMap[string(key)] = append(val, string(value))
			} else {
				paramsMap[string(key)] = []string{string(value)}
			}
		})
		for key, value := range paramsMap {
			log.Printf("key %s value %s", key, value)
			switch string(key) {
			case "page":
				if p, e := strconv.Atoi(value[0]); e == nil {
					page.Page = p
				}
			case "size":
				if s, e := strconv.Atoi(value[0]); e == nil {
					page.Size = s
				}
			case "sort":
				if err := sort.UnmarshalText([]byte(value[0])); err == nil {
					ctxFilter.Sort = &sort
				}
			default:
				f := KeyvalueOperator{}
				if err := f.UnmarshalQueryParam(key, value...); err == nil {
					dataFilters = append(dataFilters, &f)
				}
			}
		}
		ctxFilter.Filters = dataFilters
		if err := sort.Validate(); err == nil {
			ctxFilter.Sort = &sort
		}
		if err := page.Validate(); err == nil {
			ctxFilter.Page = &page
		}
		c.Locals(cfg.ContextKey, ctxFilter)
		return c.Next()
	}
}

package filter

import "github.com/gofiber/fiber/v2"

const ContextKeyDefault = "filter"

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	ContextKey string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:       nil,
	ContextKey: ContextKeyDefault,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}

	return cfg
}

package consts

// AUTHORITIES_KEY hold the key for logged in authorities in JWT token
const AUTHORITIES_KEY = "auth"

// FIBER_CONTEXT_KEY hold the context key for c *fiber.Context to subsequence action to access via c.Local()
const FIBER_CONTEXT_KEY = "user"

// JWTSECRET hold the JWT secret for encode and decode
var JWTSECRET string

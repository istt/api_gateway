package utils

import (
	"fmt"

	"math/rand"
	"time"
	"unsafe"

	"github.com/istt/api_gateway/pkg/fiber/authjwt/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GetCurrentUserLogin return the current login name for this particular user
func GetCurrentUserLogin(c *fiber.Ctx) (string, error) {
	ctx := c.Locals(consts.FIBER_CONTEXT_KEY)
	if ctx == nil {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "missing JWT information")
	}
	token, ok := ctx.(*jwt.Token)
	if !ok {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "invalid JWT information")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fiber.NewError(fiber.StatusExpectationFailed, "invalid JWT information")
	}
	return fmt.Sprint(claims["sub"]), nil
}

// RandStringBytesMaskImprSrcUnsafe generate random string
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

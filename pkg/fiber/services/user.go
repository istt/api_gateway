package services

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

var USER_SVC UserService

type UserService interface {
	// CheckPasswordHash compare password with hash
	CheckPasswordHash(password, hash string) bool

	// GetUserByUsername return user information based on login username
	GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error)

	// HashPassword hash the given password with bcrypt method
	HashPassword(password string) (string, error)

	// IsValidToken check if current login match the given jwt subject
	IsValidToken(t *jwt.Token, login string) bool

	// RegisterAccount register for a new account
	RegisterAccount(ctx context.Context, account *shared.ManagedUserDTO) error

	// Save save the account info for existing account
	SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error
}

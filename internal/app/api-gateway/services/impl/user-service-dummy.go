package impl

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

// UserServiceDummy act as a placeholders for demo purpose of how to create the implementation for this service
type UserServiceDummy struct {
}

// NewUserServiceDummy create the single-ton instance for this service
func NewUserServiceDummy() services.UserService {
	return &UserServiceDummy{}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceDummy) CheckPasswordHash(password, hash string) bool {
	return false
}

// GetUserByUsername return user information based on login username
func (svc *UserServiceDummy) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	return nil, fmt.Errorf("not implemented")
}

// HashPassword hash the given password with bcrypt method
func (svc *UserServiceDummy) HashPassword(password string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceDummy) IsValidToken(t *jwt.Token, login string) bool {
	return false
}

// IsValidUser validate one user retrieve from etcd
func (svc *UserServiceDummy) IsValidUser(ctx context.Context, login string, password string) bool {
	return false
}

// RegisterAccount register for a new account
func (svc *UserServiceDummy) RegisterAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// SaveAccount save the current account
func (svc *UserServiceDummy) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

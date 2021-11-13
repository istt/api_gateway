package impl

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceDefault use bcrypt to hash password, and inject userRepository to manage user data.
// It go in pair with the `tools/account-cli` application to manage user from command line.
type UserServiceDefault struct {
	r services.UserRepository
}

// NewUserServiceDefault create the single-ton instance for this service
func NewUserServiceDefault(userRepository services.UserRepository) services.UserService {
	return &UserServiceDefault{r: userRepository}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceDefault) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil

}

// GetUserByUsername return user information based on login username
func (svc *UserServiceDefault) GetUserByUsername(_ context.Context, login string) (*shared.ManagedUserDTO, error) {
	userInfo, err := svc.r.FindByLogin(login)
	return &userInfo, err
}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceDefault) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err

}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceDefault) IsValidToken(t *jwt.Token, login string) bool {
	if err := t.Claims.Valid(); err != nil {
		return false
	}
	if RegisteredClaims, ok := t.Claims.(jwt.RegisteredClaims); ok {
		return RegisteredClaims.Subject == login
	}
	if mapClaims, ok := t.Claims.(jwt.MapClaims); ok {
		if subject, ok := mapClaims["sub"]; ok {
			return subject == login
		}
	}
	return false
}

// RegisterAccount register for a new account
func (svc *UserServiceDefault) RegisterAccount(_ context.Context, account *shared.ManagedUserDTO) error {
	return svc.r.Save(account)
}

// SaveAccount save the current account
func (svc *UserServiceDefault) SaveAccount(_ context.Context, account *shared.ManagedUserDTO) error {
	return svc.r.Save(account)
}

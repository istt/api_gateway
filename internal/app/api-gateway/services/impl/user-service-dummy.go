package impl

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

// UserServiceDummy act as a placeholders for demo purpose of how to create the implementation for this service
// following username / email / password / authorities are availble:
//   admin@localhost / admin / admin / ROLE_ADMIN, ROLE_USER
//	 user@localhost / user / user / ROLE_USER
type UserServiceDummy struct {
}

// NewUserServiceDummy create the single-ton instance for this service
func NewUserServiceDummy() services.UserService {
	return &UserServiceDummy{}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceDummy) CheckPasswordHash(password, hash string) bool {
	return password == hash
}

// GetUserByUsername return user information based on login username
func (svc *UserServiceDummy) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	if (login == "admin") || (login == "admin@localhost") {
		return &shared.ManagedUserDTO{
			UserDTO: shared.UserDTO{
				Id:          "admin",
				Login:       "admin",
				Email:       "admin@localhost",
				FirstName:   "Admin",
				LastName:    "Dummy",
				LangKey:     "en",
				Activated:   true,
				Authorities: []string{"ROLE_USER", "ROLE_ADMIN"},
			},
			CreatedBy:        "system",
			CreatedDate:      "2006-01-02",
			Password:         "admin",
			LastModifiedBy:   "system",
			LastModifiedDate: "2006-01-02",
		}, nil
	}
	if (login == "user") || (login == "user@localhost") {
		return &shared.ManagedUserDTO{
			UserDTO: shared.UserDTO{
				Id:          "user",
				Login:       "user",
				Email:       "user@localhost",
				FirstName:   "Dummy",
				LastName:    "User",
				LangKey:     "en",
				Activated:   true,
				Authorities: []string{"ROLE_USER"},
			},
			CreatedBy:        "system",
			CreatedDate:      "2006-01-02",
			Password:         "user",
			LastModifiedBy:   "system",
			LastModifiedDate: "2006-01-02",
		}, nil
	}
	return nil, fmt.Errorf("not implemented")
}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceDummy) HashPassword(password string) (string, error) {
	return password, nil
}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceDummy) IsValidToken(t *jwt.Token, login string) bool {
	if err := t.Claims.Valid(); err != nil {
		return false
	}
	if standardClaims, ok := t.Claims.(jwt.StandardClaims); ok {
		return standardClaims.Subject == login
	}
	if mapClaims, ok := t.Claims.(jwt.MapClaims); ok {
		if subject, ok := mapClaims["sub"]; ok {
			return subject == login
		}
	}
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

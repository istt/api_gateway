package impl

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"github.com/tidwall/buntdb"

	"golang.org/x/crypto/bcrypt"
)

// UserServiceBuntDB act as a placeholders for demo purpose of how to create the implementation for this service
// following username / email / password / authorities are availble:
//   admin@localhost / admin / admin / ROLE_ADMIN, ROLE_USER
//	 user@localhost / user / user / ROLE_USER
type UserServiceBuntDB struct {
	prefix string
}

// NewUserServiceBuntDB create the single-ton instance for this service
func NewUserServiceBuntDB(prefix string) services.UserService {
	app.BuntDB.CreateIndex(prefix, prefix+"*")
	return &UserServiceBuntDB{
		prefix: prefix,
	}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceBuntDB) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceBuntDB) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GetUserByUsername return user information based on login username
func (svc *UserServiceBuntDB) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	userInfo := &shared.ManagedUserDTO{}
	err := app.BuntDB.View(func(tx *buntdb.Tx) error {
		userdata, err := tx.Get(svc.prefix + login)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(userdata), userInfo)

	})
	return userInfo, err
}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceBuntDB) IsValidToken(t *jwt.Token, login string) bool {
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
func (svc *UserServiceBuntDB) RegisterAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return app.BuntDB.Update(func(tx *buntdb.Tx) error {
		userdata, err := json.Marshal(account)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(svc.prefix+account.Login, string(userdata), nil)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(svc.prefix+account.Email, string(userdata), nil)
		if err != nil {
			return err
		}
		return nil
	})

}

// SaveAccount save the current account
func (svc *UserServiceBuntDB) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	return app.BuntDB.Update(func(tx *buntdb.Tx) error {
		userdata, err := json.Marshal(account)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(svc.prefix+account.Login, string(userdata), nil)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(svc.prefix+account.Email, string(userdata), nil)
		if err != nil {
			return err
		}
		return nil
	})
}

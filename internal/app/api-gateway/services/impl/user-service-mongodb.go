package impl

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/form3tech-oss/jwt-go"
	"github.com/istt/api_gateway/internal/app/api-gateway/instances"
	"github.com/istt/api_gateway/pkg/fiber/middleware"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceMongodb act as a placeholders for demospurpose of how to create the implementation for this service
// following username / email / password / authorities are availble:
//   admin@localhost / admin / admin / ROLE_ADMIN, ROLE_USER
//	 user@localhost / user / user / ROLE_USER
type UserServiceMongodb struct {
	Db *mongo.Database
}

// NewUserServiceMongodb create the single-ton instance for this service
func NewUserServiceMongodb() services.UserService {
	return &UserServiceMongodb{}
}

// ham generate
func GenerateFromPassword(password []byte, cost int) ([]byte, error)

//hash truyen vao password dang string voi ham gre
func (svc *UserServiceMongodb) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err

}

// CheckPasswordHash compare password with hash
func (svc *UserServiceMongodb) CheckPasswordHash(password, hash string) bool {
	return password == hash
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GetUserByUsername return user information based on login username
func (svc *UserServiceMongodb) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	if (login == "admin") || (login == "admin@localhost") {
		return &shared.ManagedUserDTO{
			UserDTO: shared.UserDTO{
				Id:          "admin",
				Login:       "admin",
				Email:       "admin@localhost",
				FirstName:   "Admin",
				LastName:    "Mongodb",
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
				FirstName:   "Mongodb",
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
func (svc *UserServiceMongodb) HashPassword(password string) (string, error) {
	return password, nil
}

// IsValidToken check if current login match the given jwt subject
func (svc *UserServiceMongodb) IsValidToken(t *jwt.Token, login string) bool {
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
func (svc *UserServiceMongodb) RegisterAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	userdata, _ := bson.Marshal(account)
	_, err := mongo.Db.collection("account").InsertOne(context.Background(), userdata)
	if err != nil {
		return err
	}
	return nil

	if govalidator.IsNull(account.Login) || govalidator.IsNull(account.Email) || govalidator.IsNull(account.Password) {
		c.JSON(bson.D, 400, "Data can not empty")
		return nil
	}

	if !govalidator.IsEmail(account.Email) {
		c.JSON(bson.D, 400, "Email is invalid")
		return nil
	}

	if err != nil {
		c.JSON(bson.D, 500, "Register has failed")
		return err
	}

	newUser := bson.M{"username": account.Login, "email": account.Email, "password": account.Password}

	_, errs := collection.InsertOne(context.TODO(), newUser)

	if errs != nil {
		c.JSON(bson.D, 500, "Register has failed")
		return err
	}

	c.JSON(bson.D, 201, "Register Succesfully")
	return nil
}

// SaveAccount save the current account
func (svc *UserServiceMongodb) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	var updatedInfo shared.UserDTO
	if err := c.BodyParser(&updatedInfo); err != nil {
		return err
	}
	logininfo, err := middleware.GetCurrentUserLogin(c)
	if err != nil {
		return err
	}
	account, err = instances.UserService.GetUserByUsername(c.Context(), logininfo)
	if err != nil {
		return err
	}

	account.UserDTO = updatedInfo
	if err := instances.UserService.SaveAccount(c.Context(), account); err != nil {
		return err
	}
	return c.JSON(account.UserDTO)
}

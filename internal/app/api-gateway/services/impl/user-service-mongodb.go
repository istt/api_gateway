package impl

import (
	"context"
	"fmt"
	"log"

	"github.com/form3tech-oss/jwt-go"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceMongodb act as a placeholders for demospurpose of how to create the implementation for this service
// following username / email / password / authorities are availble:
//   admin@localhost / admin / admin / ROLE_ADMIN, ROLE_USER
//	 user@localhost / user / user / ROLE_USER
type UserServiceMongodb struct {
	userCollection *mongo.Collection
}

// NewUserServiceMongodb create the single-ton instance for this service
func NewUserServiceMongodb() services.UserService {
	return &UserServiceMongodb{
		userCollection: app.MongoDB.Collection("user"),
	}
}

// CheckPasswordHash compare password with hash
func (svc *UserServiceMongodb) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

// GetUserByUsername return user information based on login username
func (svc *UserServiceMongodb) GetUserByUsername(ctx context.Context, login string) (*shared.ManagedUserDTO, error) {
	result := &shared.ManagedUserDTO{}
	// Search mongodb collection for user info
	if err := svc.userCollection.FindOne(ctx, bson.M{"login": login}).Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result, nil

}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceMongodb) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println(err)

	}
	return string(hash), nil
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

	Insertdata, err := svc.userCollection.InsertOne(ctx, account)
	if err != nil {
		fmt.Printf("Insert info error %s", err)
		return err
	} else {
		fmt.Printf("Insert info sucess!! %+v", Insertdata)
		return nil
	}
}

// SaveAccount save the current account
func (svc *UserServiceMongodb) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	if account.Id == "" {
		res, err := svc.userCollection.InsertOne(context.Background(), account)
		if err != nil {
			return err
		}
		if objId, ok := res.InsertedID.(primitive.ObjectID); ok {
			account.Id = objId.Hex()
		}
	} else {
		res, err := svc.userCollection.UpdateByID(context.Background(), account.Id, account)
		if err != nil {
			return err
		}
		if res.MatchedCount == 0 {
			return fmt.Errorf("no records match")
		}
	}
	return nil
}

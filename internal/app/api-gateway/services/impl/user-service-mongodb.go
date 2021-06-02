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
		return nil, err
	}
	return result, nil

}

// HashPassword hash the given password with any kind of encrypt for password. Can be MD5, SHA1 or BCrypt
func (svc *UserServiceMongodb) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
	Insertdata, err := svc.userCollection.InsertOne(context.Background(), account)
	if err != nil {
		return err
	}

	fmt.Println("Insert sucess!!", Insertdata.InsertedID)
	return nil
}

// SaveAccount save the current account
func (svc *UserServiceMongodb) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {
	users := svc.userCollection
	saveRes, err := users.InsertOne(context.Background(), account)
	if saveRes != nil {
		return err

	} else {
		return nil
	}

}
func (svc *UserServiceMongodb) EditInfor(ctx context.Context, account *shared.ManagedUserDTO) error {
	filter := &shared.ManagedUserDTO{}
	editor := bson.D{{
		Key:   "login",
		Value: nil,
	}}
	editResult, err := svc.userCollection.UpdateOne(context.TODO(), filter, editor)

	if err != nil {
		log.Fatal(err)

	}
	fmt.Printf("Matched %v edit %v documents.\n", editResult.MatchedCount, editResult.ModifiedCount)
	return nil
}

func (svc *UserServiceMongodb) DeleteAccountInfo(ctx context.Context, account *shared.ManagedUserDTO) error {
	Delete := &shared.ManagedUserDTO{}
	deleteResult, err := svc.userCollection.DeleteMany(context.TODO(), Delete)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v in collection.", deleteResult.DeletedCount)
	return nil
}

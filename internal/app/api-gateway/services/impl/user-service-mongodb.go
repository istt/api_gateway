package impl

import (
	"context"
	"fmt"
	"log"
	"reflect"

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

	Insertdata, err := svc.userCollection.InsertOne(ctx, account)
	if err != nil {
		fmt.Println("Insert info error", err)
		return err
	} else {
		fmt.Println("Insert info sucess!!", reflect.TypeOf(Insertdata))
		return nil
	}
}

// SaveAccount save the current account
func (svc *UserServiceMongodb) SaveAccount(ctx context.Context, account *shared.ManagedUserDTO) error {

	saveRes, err := svc.userCollection.InsertOne(context.Background(), account)
	if saveRes != nil {
		return err
	} else {
		fmt.Println("save sucess!!!", saveRes.InsertedID)
		return nil
	}

}

//change password
func (svc *UserServiceMongodb) editAccount(ctx context.Context, login string, email string, password string, firstName string, lastName string) (*shared.ManagedUserDTO, error) {
	editInfo := &shared.ManagedUserDTO{}
	result, err := svc.userCollection.UpdateByID(ctx, bson.M{"login": login},
		bson.M{"$set": bson.M{"password": password, "email": email, "firstName": firstName, "lastName": lastName}})

	if err != nil {
		log.Fatal(err)

	}
	fmt.Printf("updated %v documents.\n", result.ModifiedCount)
	return editInfo, nil
}

func (svc *UserServiceMongodb) changePassword(ctx context.Context, currentPassword string, newPassword string) (*shared.PasswordChangeDTO, error) {
	changepassword := &shared.PasswordChangeDTO{}
	changeRes, err := svc.userCollection.UpdateOne(ctx, bson.M{"currentPassword": currentPassword},
		bson.M{"$set": bson.M{"NewPassword": newPassword}})

	if err != nil {
		log.Fatal(err)
		fmt.Println("fail!!", err)
	} else {
		fmt.Println("password changed!!\n", changeRes.ModifiedCount)

	}
	return changepassword, nil
}

// delete account
func (svc *UserServiceMongodb) DeleteAccount(ctx context.Context, account *shared.ManagedUserDTO) error {

	deleteResult, err := svc.userCollection.DeleteMany(context.TODO(), account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v in collection.", deleteResult.DeletedCount)
	return nil
}

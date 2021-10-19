package repositories

import (
	"fmt"

	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"go.mongodb.org/mongo-driver/mongo"
)

const USER_PREFIX = "U:"

type UserRepositoryMongoDB struct {
	db *mongo.Collection
}

func NewUserRepositoryMongoDB(db *mongo.Collection) services.UserRepository {
	return &UserRepositoryMongoDB{
		db: db,
	}
}

// Returns the number of entities available.
func (r *UserRepositoryMongoDB) Count() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// Deletes a given entity.
func (r *UserRepositoryMongoDB) Delete(entity shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Deletes all entities managed by the repository.
func (r *UserRepositoryMongoDB) DeleteAll() error {
	return fmt.Errorf("not implemented")
}

// Deletes the given entities.
func (r *UserRepositoryMongoDB) DeleteAllEntities(entities []shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Deletes all instances of the type T with the given IDs.
func (r *UserRepositoryMongoDB) DeleteAllById(ids []string) error {
	return fmt.Errorf("not implemented")
}

// Deletes the entity with the given id.
func (r *UserRepositoryMongoDB) DeleteById(ID string) error {
	return fmt.Errorf("not implemented")
}

// Returns whether an entity with the given id exists.
func (r *UserRepositoryMongoDB) ExistsById(ID string) bool {
	return false
}

// ExistsByLogin check if user is exists by Login
func (r *UserRepositoryMongoDB) ExistsByLogin(ID string) bool {
	return false
}

// Returns all instances of the type.
func (r *UserRepositoryMongoDB) FindAll() ([]shared.UserDTO, error) {
	return []shared.UserDTO{}, fmt.Errorf("not implemented")
}

// Returns all instances of the type T with the given IDs.
func (r *UserRepositoryMongoDB) FindAllById(ids []string) ([]shared.UserDTO, error) {
	return []shared.UserDTO{}, fmt.Errorf("not implemented")
}

// Retrieves an entity by its id.
func (r *UserRepositoryMongoDB) FindById(ID string) (shared.UserDTO, error) {
	return shared.UserDTO{}, fmt.Errorf("not implemented")
}

func (r *UserRepositoryMongoDB) FindByLogin(login string) (shared.ManagedUserDTO, error) {
	return shared.ManagedUserDTO{}, fmt.Errorf("not implemented")
}

// Saves a given entity.
func (r *UserRepositoryMongoDB) Save(entity *shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Saves all given entities.
func (r *UserRepositoryMongoDB) SaveAll(entities []*shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

func (r *UserRepositoryMongoDB) FindAllAuthorities() ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

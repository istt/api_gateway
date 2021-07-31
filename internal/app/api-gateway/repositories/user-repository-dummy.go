package repositories

import (
	"fmt"

	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
)

type UserRepositoryDummy struct {
}

func NewUserRepositoryDummy() services.UserRepository {
	return &UserRepositoryDummy{}
}

// Returns the number of entities available.
func (r *UserRepositoryDummy) Count() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// Deletes a given entity.
func (r *UserRepositoryDummy) Delete(entity shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Deletes all entities managed by the repository.
func (r *UserRepositoryDummy) DeleteAll() error {
	return fmt.Errorf("not implemented")
}

// Deletes the given entities.
func (r *UserRepositoryDummy) DeleteAllEntities(entities []shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Deletes all instances of the type T with the given IDs.
func (r *UserRepositoryDummy) DeleteAllById(ids []string) error {
	return fmt.Errorf("not implemented")
}

// Deletes the entity with the given id.
func (r *UserRepositoryDummy) DeleteById(ID string) error {
	return fmt.Errorf("not implemented")
}

// Returns whether an entity with the given id exists.
func (r *UserRepositoryDummy) ExistsById(ID string) bool {
	return false
}

// ExistsByLogin check if user is exists by Login
func (r *UserRepositoryDummy) ExistsByLogin(ID string) bool {
	return false
}

// Returns all instances of the type.
func (r *UserRepositoryDummy) FindAll() ([]shared.UserDTO, error) {
	return []shared.UserDTO{}, fmt.Errorf("not implemented")
}

// Returns all instances of the type T with the given IDs.
func (r *UserRepositoryDummy) FindAllById(ids []string) ([]shared.UserDTO, error) {
	return []shared.UserDTO{}, fmt.Errorf("not implemented")
}

// Retrieves an entity by its id.
func (r *UserRepositoryDummy) FindById(ID string) (shared.UserDTO, error) {
	return shared.UserDTO{}, fmt.Errorf("not implemented")
}

func (r *UserRepositoryDummy) FindByLogin(login string) (shared.ManagedUserDTO, error) {
	return shared.ManagedUserDTO{}, fmt.Errorf("not implemented")
}

// Saves a given entity.
func (r *UserRepositoryDummy) Save(entity *shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

// Saves all given entities.
func (r *UserRepositoryDummy) SaveAll(entities []*shared.ManagedUserDTO) error {
	return fmt.Errorf("not implemented")
}

func (r *UserRepositoryDummy) FindAllAuthorities() ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

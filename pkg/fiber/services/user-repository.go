package services

import "github.com/istt/api_gateway/pkg/fiber/shared"

type UserRepository interface {

	// Returns the number of entities available.
	Count() (int, error)

	// Deletes a given entity.
	Delete(entity shared.ManagedUserDTO) error

	// Deletes all entities managed by the repository.
	DeleteAll() error

	// Deletes the given entities.
	DeleteAllEntities(entities []shared.ManagedUserDTO) error

	// Deletes all instances of the type T with the given IDs.
	DeleteAllById(ids []string) error

	// Deletes the entity with the given id.
	DeleteById(ID string) error

	// Returns whether an entity with the given id exists.
	ExistsByLogin(login string) bool

	// Returns all instances of the type.
	FindAll() ([]shared.UserDTO, error)

	// Returns all instances of the type T with the given IDs.
	FindAllById(ids []string) ([]shared.UserDTO, error)

	// Retrieves an entity by its id.
	FindById(ID string) (shared.UserDTO, error)

	// Retrieves an entity by its id.
	FindByLogin(login string) (shared.ManagedUserDTO, error)

	// Saves a given entity.
	Save(entity *shared.ManagedUserDTO) error

	// Saves all given entities.
	SaveAll(entities []*shared.ManagedUserDTO) error

	// Authorities
	FindAllAuthorities() ([]string, error)
}

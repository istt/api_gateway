package repositories

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"github.com/tidwall/buntdb"
)

type UserRepositoryBuntDB struct {
	db *buntdb.DB
}

func NewUserRepositoryBuntDB(db *buntdb.DB) services.UserRepository {
	return &UserRepositoryBuntDB{db: db}
}

// Returns the number of entities available.
func (r *UserRepositoryBuntDB) Count() (int, error) {
	cnt := 0
	r.db.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys("", func(key, value string) bool {
			cnt++
			return true
		})
		return nil
	})
	return cnt, nil
}

// Deletes a given entity.
func (r *UserRepositoryBuntDB) Delete(entity shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(entity.Id)
		return err
	})
}

// Deletes all entities managed by the repository.
func (r *UserRepositoryBuntDB) DeleteAll() error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		return tx.DeleteAll()
	})
}

// Deletes the given entities.
func (r *UserRepositoryBuntDB) DeleteAllEntities(entities []shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, e := range entities {
			if _, err := tx.Delete(e.Id); err != nil {
				return err
			}
		}
		return nil
	})
}

// Deletes all instances of the type T with the given IDs.
func (r *UserRepositoryBuntDB) DeleteAllById(ids []string) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, id := range ids {
			if _, err := tx.Delete(id); err != nil {
				return err
			}
		}
		return nil
	})
}

// Deletes the entity with the given id.
func (r *UserRepositoryBuntDB) DeleteById(ID string) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(ID)
		return err
	})
}

// Returns whether an entity with the given id exists.
func (r *UserRepositoryBuntDB) ExistsById(ID string) bool {
	return !errors.Is(r.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(ID, false)
		return err
	}), buntdb.ErrNotFound)
}

// ExistsByLogin check if user is exists by Login
func (r *UserRepositoryBuntDB) ExistsByLogin(ID string) bool {
	return r.ExistsById(ID)
}

// Returns all instances of the type.
func (r *UserRepositoryBuntDB) FindAll() ([]shared.UserDTO, error) {
	entities := make([]shared.UserDTO, 0)
	err := r.db.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys("", func(key, value string) bool {
			entity := shared.UserDTO{}
			if err := json.Unmarshal([]byte(value), &entity); err == nil {
				entities = append(entities, entity)
			}
			return true
		})
		return nil
	})

	return entities, err
}

// Returns all instances of the type T with the given IDs.
func (r *UserRepositoryBuntDB) FindAllById(ids []string) ([]shared.UserDTO, error) {
	entities := make([]shared.UserDTO, 0)
	err := r.db.View(func(tx *buntdb.Tx) error {
		for _, id := range ids {
			entity := shared.UserDTO{}
			if value, err := tx.Get(id); err == nil {
				if err := json.Unmarshal([]byte(value), &entity); err == nil {
					entities = append(entities, entity)
				}
			}
		}
		return nil
	})

	return entities, err
}

// Retrieves an entity by its id.
func (r *UserRepositoryBuntDB) FindById(ID string) (shared.UserDTO, error) {
	entity := shared.UserDTO{}
	err := r.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.Get(ID)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &entity); err != nil {
			return err
		}
		return nil
	})

	return entity, err
}

func (r *UserRepositoryBuntDB) FindByLogin(login string) (shared.ManagedUserDTO, error) {
	entity := shared.ManagedUserDTO{}
	err := r.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.Get(login)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &entity); err != nil {
			return err
		}
		return nil
	})

	return entity, err
}

// Saves a given entity.
func (r *UserRepositoryBuntDB) Save(entity *shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		val, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(entity.Id, string(val), nil)
		return err
	})
}

// Saves all given entities.
func (r *UserRepositoryBuntDB) SaveAll(entities []*shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, entity := range entities {
			val, err := json.Marshal(entity)
			if err != nil {
				return err
			}
			if _, _, err = tx.Set(entity.Id, string(val), nil); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepositoryBuntDB) FindAllAuthorities() ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}

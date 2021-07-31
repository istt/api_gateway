package repositories

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/pkg/fiber/services"
	"github.com/istt/api_gateway/pkg/fiber/shared"
	"github.com/tidwall/buntdb"
)

const USER_PREFIX = "U:"

type UserRepositoryBuntDB struct {
	db *buntdb.DB
}

func NewUserRepositoryBuntDB(dbpath string) services.UserRepository {
	db, err := buntdb.Open(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	db.CreateIndex(USER_PREFIX, USER_PREFIX+"*", buntdb.IndexString)
	return &UserRepositoryBuntDB{
		db: db,
	}
}

// Returns the number of entities available.
func (r *UserRepositoryBuntDB) Count() (int, error) {
	cnt := 0
	err := r.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(USER_PREFIX, func(key, value string) bool {
			cnt++
			return true
		})
	})
	return cnt, err
}

// Deletes a given entity.
func (r *UserRepositoryBuntDB) Delete(entity shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(USER_PREFIX + entity.Id)
		if err != nil {
			return err
		}
		return nil
	})
}

// Deletes all entities managed by the repository.
func (r *UserRepositoryBuntDB) DeleteAll() error {
	ids := make([]string, 0)
	err := r.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(USER_PREFIX, func(key, value string) bool {
			ids = append(ids, key)
			return true
		})
	})
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, id := range ids {
			_, err := tx.Delete(id)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Deletes the given entities.
func (r *UserRepositoryBuntDB) DeleteAllEntities(entities []shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, entity := range entities {
			_, err := tx.Delete(entity.Id)
			if err != nil {
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
			_, err := tx.Delete(USER_PREFIX + id)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Deletes the entity with the given id.
func (r *UserRepositoryBuntDB) DeleteById(ID string) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(USER_PREFIX + ID)
		return err
	})
}

// Returns whether an entity with the given id exists.
func (r *UserRepositoryBuntDB) ExistsById(ID string) bool {
	return r.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(USER_PREFIX + ID)
		return err
	}) == nil
}

func (r *UserRepositoryBuntDB) ExistsByLogin(ID string) bool {
	return r.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(USER_PREFIX + ID)
		return err
	}) == nil
}

// Returns all instances of the type.
func (r *UserRepositoryBuntDB) FindAll() ([]shared.UserDTO, error) {
	entities := make([]shared.UserDTO, 0)
	err := r.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(USER_PREFIX, func(key, value string) bool {
			var entity shared.ManagedUserDTO
			if err := json.Unmarshal([]byte(value), &entity); err != nil {
				log.Printf("error unmashal data: %s", err)
				return false
			}
			entities = append(entities, entity.UserDTO)
			return true
		})
	})
	return entities, err
}

// Returns all instances of the type T with the given IDs.
func (r *UserRepositoryBuntDB) FindAllById(ids []string) ([]shared.UserDTO, error) {
	entities := make([]shared.UserDTO, 0)
	err := r.db.View(func(tx *buntdb.Tx) error {
		for _, id := range ids {
			value, err := tx.Get(USER_PREFIX + id)
			if err != nil {
				continue
			}
			var entity shared.ManagedUserDTO
			if err := json.Unmarshal([]byte(value), &entity); err == nil {
				entities = append(entities, entity.UserDTO)
			}
		}
		return nil
	})
	return entities, err
}

// Retrieves an entity by its id.
func (r *UserRepositoryBuntDB) FindById(ID string) (shared.UserDTO, error) {
	var entity shared.ManagedUserDTO
	err := r.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.Get(USER_PREFIX + ID)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &entity); err != nil {
			return err
		}
		return nil
	})
	return entity.UserDTO, err
}

func (r *UserRepositoryBuntDB) FindByLogin(login string) (shared.ManagedUserDTO, error) {
	var entity shared.ManagedUserDTO
	err := r.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.Get(USER_PREFIX + login)
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
	if entity.Id == "" {
		entity.Id = entity.Login
	}
	return r.db.Update(func(tx *buntdb.Tx) error {
		entityJson, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(USER_PREFIX+entity.Id, string(entityJson), nil)
		return err
	})
}

// Saves all given entities.
func (r *UserRepositoryBuntDB) SaveAll(entities []*shared.ManagedUserDTO) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		for _, entity := range entities {
			entityJson, err := json.Marshal(entity)
			if err != nil {
				return err
			}
			_, _, err = tx.Set(USER_PREFIX+entity.Id, string(entityJson), nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepositoryBuntDB) FindAllAuthorities() ([]string, error) {
	authorities := app.Config.Strings("security.authorities")
	if len(authorities) == 0 {
		return authorities, fmt.Errorf("invalid authorities")
	}
	return authorities, nil
}

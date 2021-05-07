package test

import (
	"testing"

	"github.com/istt/api_gateway/internal/app"
	"github.com/tidwall/buntdb"
)

func TestBuntDB(t *testing.T) {
	app.BuntDBConfig()
	app.BuntDBInit()

	// + buntdb update key
	app.BuntDB.Update(func(tx *buntdb.Tx) error {
		previousValue, replaced, err := tx.Set("testKey", "testValue", nil)
		t.Logf("previous %s replaced %v", previousValue, replaced)
		return err
	})

	// + buntdb look-up key
	app.BuntDB.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			t.Logf("%s -> %s", key, value)
			return true
		})
	})

}

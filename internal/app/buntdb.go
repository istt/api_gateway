package app

import (
	"log"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/tidwall/buntdb"
)

var (
	// DBConn hold the connection to database
	BuntDB *buntdb.DB
)

// BuntDBConfig configure application runtime
func BuntDBConfig() {
	// koanf defautl values
	Config.Load(confmap.Provider(map[string]interface{}{
		"buntdb.path": ":memory:",
	}, "."), nil)
}

// BuntDBInit initiate database
func BuntDBInit() {
	var err error
	log.Printf("Connecting to database: %s", Config.String("buntdb.path"))
	BuntDB, err = buntdb.Open(Config.String("buntdb.path"))
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connection Opened to BuntDB")
}

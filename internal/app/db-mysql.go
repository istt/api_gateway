package app

import (
	"fmt"
	"log"

	"github.com/knadh/koanf/providers/confmap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DBConn hold the connection to database
	MysqlDB *gorm.DB
)

// MysqlDBConfig configure application runtime
// link: https://gorm.io/docs/connecting_to_the_database.html
func MysqlDBConfig() {
	// koanf defautl values
	Config.Load(confmap.Provider(map[string]interface{}{
		// "mysql.dsn": "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local",
		"mysql.user": "root",
		"mysql.pass": "",
		"mysql.host": "127.0.0.1",
		"mysql.port": 3306,
		"mysql.name": "api_gateway",
	}, "."), nil)
}

// MysqlDBInit initiate database
func MysqlDBInit() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		Config.MustString("mysql.user"),
		Config.String("mysql.pass"),
		Config.MustString("mysql.host"),
		Config.MustInt("mysql.port"),
		Config.MustString("mysql.name"),
	)
	log.Printf("Connecting to database: %s", Config.String("db.dsn"))
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	MysqlDB = db
	log.Println("Connection Opened to MysqlDB")
}

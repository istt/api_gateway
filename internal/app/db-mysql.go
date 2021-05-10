package app

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBMysqlConfig() {
	//Open connection to database
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:                "MysqlDriver",
		DSN:                       "",
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  new(int),
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: false,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	defer db.Close(gorm.Config)
}

var (
	Client *sql.DB

	username = os.Getenv("MYSQL_IDS_USERNAME")
	password = os.Getenv("MYSQL_IDS_PASSWORD")
	host     = os.Getenv("MYSQL_IDS_HOST")
	port     = os.Getenv("MYSQL_IDS_PORT")
	schema   = os.Getenv("MYSQL_IDS_SCHEMA")
)

func DBMysqlInit() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		username, password, host, port, schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Error("connecting to database failed: ", err)
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)

	}
	logger.Info("database successfully configured")
}

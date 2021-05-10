package test

import (
	"testing"

	"github.com/istt/api_gateway/internal/app"
)

func Test_Mysql(t *testing.T) {
	app.MysqlDBConfig()
	app.MysqlDBInit()

	// run raw query
	rows := make([]string, 0)
	dberr := app.MysqlDB.Debug().Raw("SHOW TABLES").Scan(&rows)
	if dberr.Error != nil {
		t.Fatal((dberr.Error))
	}
	t.Logf("tables: %+v", rows)
}

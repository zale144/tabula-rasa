package util

import (
	. "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *DB

func GetDBConnection()  {
	var err error
	Db, err = Open("mysql", "zale144:pastazazube@tcp(127.0.0.1:3306)/superuser")
	if err != nil {
		print(err)
	}
}

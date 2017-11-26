package util

import (
	. "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *DB

func GetDBConnection()  {
	var err error
	Db, err = Open("mysql", "root:root@tcp(127.0.0.1:3306)/superhero")
	if err != nil {
		print(err)
	}
}

package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"database/sql"
)

var Db *sql.DB

func GetDBConnection()  {
	var err error

	Db, err = sql.Open("mysql", "zale144:pastazazube@tcp(127.0.0.1:3306)/superuser")
	if err != nil {
		log.Fatalln(err)
	}
}

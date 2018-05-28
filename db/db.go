package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var Db *sqlx.DB

func GetDBConnection()  {
	var err error

	Db, err = sqlx.Connect("mysql", "zale144:pastazazube@tcp(127.0.0.1:3306)/superuser")
	if err != nil {
		log.Fatalln(err)
	}
}

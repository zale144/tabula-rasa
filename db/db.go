package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

// read request
type ReadReq struct{
	Key string
	Resp chan *sql.DB
}
// write request
type WriteReq struct{
	Key string
	Val *sql.DB
	Resp chan bool
}
// remove request
type RemoveReq struct{
	Key string
	Resp chan bool
}
// clear request
type ClearReq struct{
	Resp chan bool
}
// read, write, remove, clear request channels
var (
	reads   = make(chan *ReadReq)
	writes  = make(chan *WriteReq)
	removes = make(chan *RemoveReq)
	clears  = make(chan *ClearReq)
)
// manage multiple database connections with a goroutine
func DBConnections()  {
	dbConnections := make(map[string]*sql.DB)
	for {
		select {
		case read := <-reads:
			read.Resp <- dbConnections[read.Key]
		case write := <-writes:
			dbConnections[write.Key] = write.Val
			write.Resp <- true
		case remove := <-removes:
			delete(dbConnections, remove.Key)    // delete value for key from map
			_, resp := dbConnections[remove.Key] // check if key still found in map
			remove.Resp <- !resp
		case clear := <-clears: // clear operation requested
			dbConnections = make(map[string]*sql.DB) // re-instantiate the  map
			clear.Resp <- len(dbConnections) == 0 // send true if map is empty now
		}
	}
}
// method for creating a database
func CreateDatabase(dbName string) error {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON `" + dbName + "`.* TO 'zale144'@'%' WITH GRANT OPTION")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = db.Exec("USE " + dbName)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// TODO fix the App.jsx so that his is not necessary
	_, err = db.Exec("CREATE TABLE example ( id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id), UNIQUE INDEX id_UNIQUE (id ASC), data varchar(32) )")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// connecting to a database by provided name
func ConnectDB(dbName string) error {
	if GetConnection(dbName) != nil {
		return nil
	}
	connection, err := sql.Open("mysql", "zale144:pastazazube@tcp(127.0.0.1:3306)/" + dbName)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Println("database '" + dbName + "' connected")
	SaveConnection(connection, dbName)
	return nil
}
// disconnecting from a database by provided name
func Disconnect(dbName string)  {
	// get the connection from the map
	connection := GetConnection(dbName)
	connection.Close()
	// remove the database connection from the map
	Remove(dbName)
}

// get connection by name
func GetConnection(name string) *sql.DB { // return db connection
	read := &ReadReq{ // instantiate a read request
		Key: name,
		Resp: make(chan *sql.DB),
	}
	reads <- read // send read request to reads channel
	ret := <- read.Resp // pull return value from the channel
	return ret
}
// add connection to map
func SaveConnection(connection *sql.DB, name string)  {
	write := &WriteReq{ // instantiate a write request
		Key: name,
		Val: connection,
		Resp: make(chan bool),
	}
	writes <- write // send write request to writes channel
	<- write.Resp   // pull boolean value from the channel
	log.Printf("Written '%s' to db connection map\n", write.Key)
}
// remove item from db connection map
func Remove(name string)  {
	remove := &RemoveReq{ // instantiate a remove request
		Key: name,
		Resp: make(chan bool),
	}
	removes <- remove // send remove request to removes channel
	<- remove.Resp // pull boolean value from the channel
	log.Printf("Removed '%s' from db connection map\n", remove.Key)
}
// clear the db connection map
func Clear()  {
	clear := &ClearReq{ // instantiate a clear request
		Resp: make(chan bool),
	}
	clears <- clear // send clear request to clears channel
	<- clear.Resp // pull boolean value from the channel
	log.Println("db connection map cleared")
}
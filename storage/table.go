package storage

import (
	"log"
	"tabula-rasa/db"
	"tabula-rasa/libs/memcache"
	"fmt"
	"encoding/json"
)

// table storage
type TableStorage struct {}

type ColumnInfo struct {
	Name      *string
	Type      *string
	Child     *interface{}
	Reference *string
}

/*------------------------------------------------------------
 * method for retrieving all entities from the given 		 *
 * database/table. It will generalize to any table structure *
 * or type of data, and read any references to other tables. *
 *-----------------------------------------------------------*/
func (ts TableStorage) Get(name, id, typ, dbName string) ([]interface{}, error) {
	entities := []interface{}{}
	// attempt to read from memcache by passing the pointer to our empty slice of interfaces
	memcache.ReadFromCache(&entities, name, id, typ, dbName)
	// if found result, return it
	if len(entities) > 0 {
		return entities, nil
	}
	// generate the query according to parameters
	query := generateSelectQuery(name, id, typ, dbName)
	// get the complete column info
	columnInfo, err := ts.getColumnInfo(name, typ, id, dbName, &query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// get the database connection by it's name
	Db := db.GetConnection(dbName)
	rows, err := Db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(columnInfo) == 0 {
		// get the table column information
		columnNames, err := rows.Columns()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for i := range columnNames {
			cI := ColumnInfo{
				Name: &columnNames[i],
			}
			columnInfo = append(columnInfo, cI)
		}
	}
	// get a list of table column types
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// map the column names into a string array
	for i, v := range colTypes {
		n := v.DatabaseTypeName()
		columnInfo[i].Type = &n
	}
	// we need to pass the pointers to the rows.Scan method
	vals := make([]interface{}, len(columnInfo))
	valPtrs := make([]interface{}, len(columnInfo))
	for i := 0; i < len(columnInfo); i++ {
		valPtrs[i] = &vals[i]
	}
	// iterate through the rows
	for rows.Next() {
		m := map[string]interface{}{}
		err = rows.Scan(valPtrs...)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		// add the retrieved values to a map
		ts.addValuesToRowMap(&vals, &m, &columnInfo, dbName)
		// convert the row map into struct interface
		e := StructBuilder{}.convertRowMapToStruct(columnInfo, m)
		// append the instance of the filled struct onto the return value slice
		entities = append(entities, e)
	}
	// if the return value was not empty, write it to cache with given parameters
	if len(entities) > 0 {
		memcache.WriteToCache(&entities, name, id, typ, dbName)
	}
	return entities, err
}
/*------------------------------------------------------------
 * method for updating/adding rows, columns, tables...	 *
 * It will generalize to any table structure or type of data *
 * and persist references to other tables.		 			 *
 *-----------------------------------------------------------*/
func (ts TableStorage) Save(name, typ, dbName string, body []byte) ([]interface{}, error) {
	var entity map[string]interface{}
	json.Unmarshal(body, &entity)

	var err error
	// generate the query string for saving the resource
	queryStr := generateSaveQueryString(name, typ, entity)
	// if columns are updated, we don't need to return the old name
	if typ == "cols" {
		delete(entity, "Old_column_name")
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(queryStr)
	// get the database connection by it's name
	Db := db.GetConnection(dbName)
	stmt, err := Db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id := entity["Id"]

	if id == "" {
		id, err = res.LastInsertId()
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		entity["Id"] = id
	} else if id == nil {
		id = entity["Table"]
	}
	// TODO only update the exact resource
	memcache.ClearCache()

	entities := []interface{}{}
	entities = append(entities, entity)
	return entities, nil
}
/*------------------------------------------------------------
 * method for deleting rows, columns, tables...	 		 *
 * It will generalize to any table structure or type of data *
 *-----------------------------------------------------------*/
func (ts TableStorage) Delete(name, id, typ, dbName string) error {
	var err error
	// generate the delete resource query string
	queryStr := generateDeleteQueryString(name, id, typ)

	fmt.Println(queryStr)
	// get the database connection by it's name
	Db := db.GetConnection(dbName)
	stmt, err := Db.Prepare(queryStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// TODO only update the exact resource
	memcache.ClearCache()
	return nil
}

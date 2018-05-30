package storage

import (
	"log"
	. "tabula-rasa/db"
	"tabula-rasa/libs/memcache"
	"fmt"
	"encoding/json"
	"reflect"
)

// table storage
type TableStorage struct {}

type ColumnInfo struct {
	Name string
	Reference interface{}
}

// function for retrieving all entities from the given database table
func (ts TableStorage) Get(name, id, typ string) ([]interface{}, error) {
	entities := []interface{}{}
	// attempt to read from memcache by passing the pointer to our empty slice of interfaces
	//memcache.ReadFromCache(&entities, name, id, typ)
	// if found result, return it
	if len(entities) > 0 {
		return entities, nil
	}
	// generate the query according to parameters
	query := generateSelectQuery(name, id, typ)

	columnInfo := []ColumnInfo{}

	if name != "tables" && typ == "rows" {
		// recursively call the Get method to get the column information for this table
		columns, err := ts.Get(name, "", "cols")
		if err != nil {
			return nil, err
		}
		// get the table column info
		columnInfo = make([]ColumnInfo, len(columns))
		for i := range columns {
			col := columns[i]
			elem := reflect.ValueOf(col).Elem()
			fld := elem.Field(0) // name of the column
			ref := elem.Field(2) // name reference table
			colName := fld.Interface()
			reference := ref.Interface()
			columnInfo[i].Name = colName.(string)
			columnInfo[i].Reference = reference.(string)

			// needs a join with referenced table
			if columnInfo[i].Reference != "" && columnInfo[i].Reference != nil {
				query = generateSelectQuery(name, id, typ)
			}
		}
	}
	rows, err := Db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	if len(columnInfo) == 0 {
		// get the table column information
		columnNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		for i := range columnNames {
			cI := ColumnInfo{
				Name: columnNames[i],
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
	typeNames := make([]string, len(colTypes))
	for i, v := range colTypes {
		typeNames[i] = v.DatabaseTypeName()
	}
	vals := make([]interface{}, len(columnInfo))
	// rows.Scan(...) expects pointers to our values
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
		for i, val := range vals {
			m[columnInfo[i].Name] = val
			// TODO use JOIN instead
			if columnInfo[i].Reference != "" && columnInfo[i].Reference != nil && val != nil {
				fmt.Printf("%v: %v = %v\n", columnInfo[i].Reference, columnInfo[i].Name, string(val.([]uint8)))
				reference, err := ts.Get(columnInfo[i].Reference.(string), string(val.([]uint8)), "rows")
				if err != nil {
					log.Println(err)
					continue
				}
				elem := reflect.ValueOf(reference[0]).Elem()
				refField := elem.Field(1) // name reference table
				m[columnInfo[i].Name] = []uint8(refField.Interface().(string))
				typeNames[i] = Text
			}
		}
		// convert the row map into struct interface
		e := StructBuilder{}.convertRowMapToStruct(columnInfo, typeNames, m)
		// append the instance of the filled struct onto the return value slice
		entities = append(entities, e)
	}
	// if the return value was not empty, write it to cache with given parameters
	if len(entities) > 0 {
		memcache.WriteToCache(&entities, name, id, typ)
	}
	return entities, err
}

// function for creating and updating entities
func (ts TableStorage) Save(name, typ string, body []byte) ([]interface{}, error) {
	var entity map[string]interface{}
	json.Unmarshal(body, &entity)

	var err error
	queryStr := generateSaveQueryString(name, typ, entity)

	if typ == "cols" {
		delete(entity, "Old_column_name")
	}
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	fmt.Println(queryStr)
	stmt, err := Db.Prepare(queryStr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec()
	if err != nil {
		log.Println(err.Error())
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
	memcache.ClearCache()

	entities := []interface{}{}
	entities = append(entities, entity)
	return entities, nil
}

// function for generating queries for deleting rows, tables and columns
func (ts TableStorage) Delete(name, id, typ string) error {
	var err error
	queryStr := generateDeleteQueryString(name, id, typ)

	fmt.Println(queryStr)
	stmt, err := Db.Prepare(queryStr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return err
	}
	memcache.ClearCache()
	return nil
}

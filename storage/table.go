package storage

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	. "tabula-rasa/db"
	"reflect"
	"tabula-rasa/libs/memcache"
	"encoding/json"
)

// table storage
type TableStorage struct {}

// function for creating new struct according to column names and types
// it takes a list of column names and their db type names,
// uses reflection to create a new type and assign fields, and returns
// the type interface that will be used as a struct.
// at first I have used a map, but I preferred the fields to be in the
// same order as in the database table
func makeEntity(columnNames, typeNames []string) interface{} {
	var sfs []reflect.StructField
	// iterate through all the column names
	for i := range columnNames {
		// set the name of the struct field
		sf := reflect.StructField{
			Name: strings.Title(columnNames[i]),
		}
		// map and set the type for the struct field
		mapTypes(typeNames[i], func(t interface{}) error {
			sf.Type = reflect.TypeOf(t)
			return nil
		})
		// append the field to the struct
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	so := reflect.New(st)
	return so.Interface()
}
// func to map struct field type according to db type name
// it takes a variadic number of functions to which we pass
// a literal of the appropriate type
func mapTypes(typeName string, setType ...func(interface{})error)  {
	inds := make([]int, 4)
	for i := 0; i < len(setType); i++ {
		inds[i] = i
	}
	switch typeName {
	case "VARCHAR", "TEXT":
		setType[inds[0]]("")
	case "INT", "INTEGER":
		setType[inds[1]](0)
	case "DOUBLE", "FLOAT":
		setType[inds[2]](0.1)
	case "BOOLEAN", "BOOL":
		setType[inds[3]](false)
	}
}
// function for retrieving all entities from the given database table
func (ts TableStorage) Get(name, id, spec string) ([]interface{}, error) {
	entities := []interface{}{}
	// attempt to read from memcache by passing the pointer to our empty slice of interfaces
	memcache.ReadFromCache(&entities, name, id, spec)
	// if found result, return it
	if len(entities) > 0 {
		return entities, nil
	}
	// generate the query according to parameters
	query, _ := generateSelectQuery(name, id, spec)
	rows, err := Db.Queryx(query)
	if err != nil {
		return nil, err
	}
	// get the table column names
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// get a list of table column types
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	// map the column names into a string array
	typeNames := make([]string, len(cols))
	for i, v := range colTypes {
		typeNames[i] = v.DatabaseTypeName()
	}
	// iterate through the rows
	for rows.Next() {
		// create the type that will hold the row data
		Type := makeEntity(cols, typeNames)
		sr := reflect.ValueOf(Type)
		e := map[string]interface{}{}
		// scan the row and write it into a map
		err := rows.MapScan(e)
		if err != nil {
			return entities, err
		}
		// iterate through the column names
		for i, c := range cols {
			if e[c] == nil {
				continue
			}
			// instantiate an array of functions that will serve as handlers for
			// setting struct fields with their appropriate data types
			funcs := []func(interface{})error{
				// handles string data types
				func(in interface{}) error {
					strVal := string(e[c].([]uint8))
					sr.Elem().Field(i).SetString(strVal)
					return nil
				},
				// handles integer data types
				func(in interface{}) error {
					intVal, err := strconv.ParseInt(string(e[c].([]uint8)), 0, 64)
					if err != nil {
						log.Fatal(err)
						return err
					}
					sr.Elem().Field(i).SetInt(intVal)
					return nil
				},
				// handles float data types
				func(in interface{}) error {
					dblVal, err := strconv.ParseFloat(string(e[c].([]uint8)),64)
					if err != nil {
						log.Fatal(err)
						return err
					}
					sr.Elem().Field(i).SetFloat(dblVal)
					return nil
				},
				// TODO add handlers for other types
			}
			// map the struct fields according to their type names
			mapTypes(typeNames[i], funcs...)
		}
		// append the instance of the filled struct onto the return value slice
		entities = append(entities, sr.Interface())
	}
	// if the return value was not empty, write it to cache with given parameters
	if len(entities) > 0 {
		memcache.WriteToCache(&entities, name, id, spec)
	}
	return entities, err
}

/*// function for generating JSON-like strings from retrieved rows from the database
func stringifyRows(name string, rows *sql.Rows, single bool) (string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	var results string
	var objs []string
	dest := make([]interface{}, len(cols))

	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return "", err
		}
		for i, raw := range rawResult {
			value := string(raw)
			if strings.HasSuffix(cols[i], "_fk") {
				childName, err := Get(name, "", "cols", cols[i])
				childName = childName[29 : len(childName)-3]
				value, err = Get(childName, "", "", value)
				if err != nil {
					return "", err
				}
			} else {
				value = `"` + value + `"`
			}
			result[i] = `"` + cols[i] + `" : ` + value
		}
		results = "{" + strings.Join(result, ", ") + "}"
		objs = append(objs, results)
	}
	var final string
	if single {
		final = strings.Join(objs, "")
	} else {
		final = "[" + strings.Join(objs, ", ") + "]"
	}
	return final, err
}*/

// function for generating SELECT queries
func generateSelectQuery(name, id, spec string) (string, bool) {
	query := "SELECT * FROM " + name
	single := false
	if id != "" && id != "0" {
		query += " WHERE id = " + id
		single = true
	} else if name == "tables" {
		query = `SELECT DISTINCT table_name FROM information_schema.columns
		WHERE column_name in ('id') AND table_schema = 'superuser';`
	} else if spec == "cols" {
		query = generateColsInfoQuery(name, "")
	}
	log.Println(query)
	return query, single
}

// function for generating query for retrieving all column information for a given table and column name
func generateColsInfoQuery(tableName, colName string) string {
	if colName != "" {
		return `select kcu.referenced_table_name
				from information_schema.key_column_usage kcu
				WHERE kcu.TABLE_SCHEMA   = 'superuser'
					AND  kcu.column_name   = '` + colName + `'
					AND  kcu.table_name = '` + tableName + `'`
	} else {
		return `SELECT cols.column_name, cols.column_type ,
				(select kcu.referenced_table_name from information_schema.key_column_usage kcu
		WHERE kcu.TABLE_SCHEMA   = 'superuser'
			AND  kcu.column_name   = cols.column_name
			AND  kcu.table_name = '` + tableName + `' limit 1) as referenced_table_name
		FROM information_schema.columns cols
			WHERE cols.TABLE_SCHEMA = 'superuser' and cols.table_name = '` + tableName + `';`
	}
}

// function for creating and updating entities
func (ts TableStorage) Save(name, typ string, body []byte) ([]interface{}, error) {
	var entity map[string]interface{}
	json.Unmarshal(body, &entity)

	var err error
	queryStr := generateInsertUpdateQuery(name, entity)
	if name == "tables" {
		queryStr = generateCreateTableQuery(entity["Table"].(string))
	}
	if typ == "cols" {
		queryStr = generateChangeColumnNameQuery(name, entity)
		delete(entity, "Old_column_name")
		memcache.RemoveFromCache(name, "", "cols")
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
		memcache.RemoveFromCache(name, "", "tabs")
	}
	memcache.RemoveFromCache(name, "", "rows")

	entities := []interface{}{}
	entities = append(entities, entity)
	return entities, nil
}

// function for changing column name
func generateChangeColumnNameQuery(tableName string, entity map[string]interface{}) string {
	queryStr := ""
	if entity["Old_column_name"] != "" && entity["Old_column_name"] != nil {
		queryStr = "ALTER TABLE " + tableName + " CHANGE `" + entity["Old_column_name"].(string) +
			"` `" + entity["Column_name"].(string) + "` " + entity["Column_type"].(string) + ";"
	} else {
		queryStr = "ALTER TABLE " + tableName + " ADD " + entity["Column_name"].(string) +
			" " + entity["Column_type"].(string) + ";"
	}
	return queryStr
}

// function for generating INSERT and UPDATE query string
func generateInsertUpdateQuery(name string, entity map[string]interface{}) string {
	queryStr := ""
	if entity["Id"] != nil && entity["Id"] != "" {
		queryStr = "UPDATE `" + name + "` SET "
		for k, v := range entity {
			if k == "Id" {
				continue
			}
			val := v.(string)
			if val == "" {
				val = " NULL ,"
			} else {
				val = "'" + v.(string) + "',"
			}
			queryStr += "`" + k + "` = " + val
		}
		queryStr = queryStr[:len(queryStr)-1] + ` WHERE id = ` + strconv.Itoa(int(entity["Id"].(float64)))
	} else {
		queryStr = "INSERT INTO `" + name + "` ( "
		columns := ""
		values := ""
		for k, v := range entity {
			if k == "Id" {
				continue
			}
			columns += "`" + k + "`" + ","
			val := v.(string)
			if val == "" {
				values += " NULL ,"
			} else {
				values += "'" + v.(string) + "',"
			}
		}
		queryStr += columns[:len(columns)-1] + " )"
		queryStr += " VALUES ( " + values[:len(values)-1] + ");"
	}
	return queryStr
}

// function for generating CREATE TABLE query string
func generateCreateTableQuery(name string) string {
	return "CREATE TABLE `" + name +
		"` ( id INT NOT NULL AUTO_INCREMENT, " +
		` PRIMARY KEY (id), UNIQUE INDEX id_UNIQUE (id ASC));`
			/*if strings.HasSuffix(val, "_id") {
				colsStr += "`" + col + "_fk` INT, "
				fk := val[:len(val)-3]
				fkStr += `,` + makeForeignKey(tableName, col, fk)
			} else {
				colsStr += " `" + col + "` " + val + ", "
			}*/
}

// function for generating the foreign key reference
func makeForeignKey(this, name, ref string) string {
	return `KEY fk_` + this + `_` + name + `_id_idx (` + name + `_fk),
			  CONSTRAINT fk_` + this + `_` + name + `_id
			  FOREIGN KEY (` + name + `_fk)
			  REFERENCES ` + ref + ` (id) ON DELETE NO ACTION ON UPDATE NO ACTION`
}

// function for generating queries for deleting rows, tables and columns
func (ts TableStorage) Delete(name, typ, id string) error {
	var err error
	var queryStr string

	if name == "tables" {
		queryStr = "DROP TABLE `" + id + "`;"
		memcache.RemoveFromCache(name, "", "tabs")
	} else if typ == "cols" {
		queryStr = "ALTER TABLE `" + name + "` DROP COLUMN `" + id + "`;"
		memcache.RemoveFromCache(name, "", "cols")
	} else {
		queryStr = "DELETE FROM " + name + " WHERE id = " + id + ";"
	}
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
	memcache.RemoveFromCache(name, "", "rows")
	return nil
}







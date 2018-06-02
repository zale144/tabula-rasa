package storage

import (
	"strconv"
	"fmt"
	"reflect"
)


// function for generating SELECT queries
func generateSelectQuery(name, id, typ string) string {
	query := "SELECT * FROM " + name
	if id != "" && id != "0" {
		query += " WHERE id = " + id
	} else if name == "tables" {
		query = `SELECT DISTINCT table_name FROM information_schema.columns
				 WHERE column_name IN ('id') AND table_schema = 'superuser';`
	} else if typ == "cols" {
		query = generateColsInfoQuery(name, "")
	}
	fmt.Println(query)
	return query
}

// function for generating query for retrieving all column information for a given table and column name
func generateColsInfoQuery(tableName, colName string) string {
	if colName != "" {
		return `SELECT kcu.referenced_table_name
				FROM information_schema.key_column_usage kcu
				WHERE kcu.TABLE_SCHEMA   = 'superuser'
				AND  kcu.column_name   = '` + colName + `'
				AND  kcu.table_name = '` + tableName + `'`
	} else {
		return `SELECT cols.column_name, cols.column_type ,
				(SELECT kcu.referenced_table_name from information_schema.key_column_usage kcu
				WHERE kcu.TABLE_SCHEMA   = 'superuser'
				AND  kcu.column_name   = cols.column_name
				AND  kcu.table_name = '` + tableName + `' limit 1) AS referenced_table_name
				FROM information_schema.columns cols
				WHERE cols.TABLE_SCHEMA = 'superuser' AND cols.table_name = '` + tableName + `';`
	}
}
// function for generating universal save query string
func generateSaveQueryString(name, typ string, entity map[string]interface{}) string {
	queryStr := generateInsertUpdateQuery(name, entity)
	if name == "tables" {
		queryStr = generateCreateTableQuery(entity["Table"].(string))
	}
	if typ == "cols" {
		queryStr = generateChangeColumnNameQuery(name, entity)
	}
	return queryStr
}
// function for generating CREATE TABLE query string
func generateCreateTableQuery(name string) string {
	return "CREATE TABLE `" + name +
		"`( id INT NOT NULL AUTO_INCREMENT, " +
		`PRIMARY KEY (id), 
		UNIQUE INDEX id_UNIQUE (id ASC));`
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
			var val string
			if reflect.TypeOf(v).String() == "string" {
				val = fmt.Sprintf("'%v',", v.(string))
			} else if reflect.TypeOf(v).String() == "float64" {
				val = fmt.Sprintf("'%v',", v.(float64))
			} else if reflect.TypeOf(v).String() == "int" {
				val = fmt.Sprintf("'%v',", v.(int))
			}
			if val == "''" {
				val = " NULL ,"
			}
			queryStr += "`" + k + "` = " + string(val)
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
	if entity["Referenced_table_name"] != "" && entity["Referenced_table_name"] != nil {
		queryStr = makeForeignKey(tableName, entity["Referenced_table_name"].(string), entity["Column_name"].(string))
	}
	return queryStr
}
// function for generating the foreign key reference
func makeForeignKey(this, name, ref string) string {
	return  ` ALTER TABLE ` + this +
			` ADD CONSTRAINT fk_` + this + `_` + ref + `_` + name +
			` FOREIGN KEY (` + ref + `)
			  REFERENCES ` + name + ` (id) ON DELETE NO ACTION ON UPDATE NO ACTION;`
}
// function for generating delete query string
func generateDeleteQueryString(name, id, typ string) string {
	queryStr := "DELETE FROM " + name + " WHERE id = " + id + ";"
	if name == "tables" {
		queryStr = "DROP TABLE `" + id + "`;"
	} else if typ == "cols" {
		queryStr = "ALTER TABLE `" + name + "` DROP COLUMN `" + id + "`;"
	}
	return queryStr
}
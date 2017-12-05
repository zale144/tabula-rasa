package dao

import (
	. "tabula-rasa/util"
	"strings"
	"fmt"
	"strconv"
	"database/sql"
)

// get all entities
func Get(args ...string) (string, error) {
	single := false
	query := generateSelectString(&single, args...)
	rows, err := Db.Query(query)
	CheckError(err)
	return stringifyRows(args[0], rows, single)
}

func stringifyRows(name string, rows *sql.Rows, single bool) (string, error) {
	cols, err := rows.Columns()
	CheckError(err)
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
		CheckError(err)
		for i, raw := range rawResult {
			value := string(raw)
			if strings.HasSuffix(cols[i], "_fk") {
				childName, err := Get(name, "", "cols", cols[i])
				childName = childName[29:len(childName)-3]
				value, err = Get(childName, value)
				CheckError(err)
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
}

func generateSelectString(single *bool, args ...string, ) string {
	mapParameters := func (args []string) (string, string, string, string) {
		params := [4]string{}
		for i, arg := range args {
			params[i] = arg
		}
		return params[0], params[1], params[2], params[3]
	}
	name, id, spec, col := mapParameters(args)
	*single = false
	query := "SELECT * FROM " + name
	if id != "" && id != "0" {
		id = args[1]
		query += " WHERE id = " + id
		*single = true
	} else if name == "home" {
		query = `SELECT DISTINCT table_name FROM information_schema.columns
		WHERE column_name in ('id') AND table_schema = 'superhero';`
	} else if spec == "cols" {
		query = makeColsInfo(name, col)
	}
	fmt.Println(query)
	return query
}

// create entity
func Create(name string, obj []byte) (string, error) {
	pairs := unmarshalJson(obj)
	queryStr, values := generateCreateString(name, pairs)
	stmt, err := Db.Prepare(queryStr)
	fmt.Println(queryStr)
	CheckError(err)
	defer stmt.Close()
	res, err := stmt.Exec(values...)
	CheckError(err)
	id, _ := res.LastInsertId()
	return Get(name, fmt.Sprint(id), "")
}

func unmarshalJson(obj []byte) [][]interface{} {
	pairList := [][]interface{}{}
	str := string(obj[2:len(obj)-2])
	pairs := strings.Split(str, `","`)
	for _, item := range pairs {
		pair := strings.Split(item, `":"`)
		couple := []interface{}{}
		for _, col := range pair {
			couple = append(couple, col)
		}
		pairList = append(pairList, couple)
	}
	return pairList
}

func generateCreateString(name string, pairs [][]interface{}) (string, []interface{}) {
	if name == "home" {
		return generateCreateTableString(pairs), []interface{}{}
	}

	values := []interface{}{}
	var id int64
	var err error
	for i := range pairs {
		if pairs[i][0] == "id" {
			id, err = strconv.ParseInt(pairs[i][1].(string), 0, 64)
			CheckError(err)
			pairs = append(pairs[:i], pairs[i+1:]...)
			break
		}
	}
	queryStr := ""
	if id > 0 {
		queryStr = "UPDATE `" + name + "` SET "
		for _, el := range pairs {
			value := strings.Replace(el[1].(string), "'", "''", -1)
			queryStr += "`" + el[0].(string) + "` = '" + value + "',"
		}
		queryStr = queryStr[:len(queryStr) - 1] + ` WHERE id = ` + strconv.FormatInt(id, 10)
		pairs = pairs[:0]
	} else {
		placeholders := ""
		for i := 0; i < len(pairs); i++ {
			placeholders += ",?"
		}
		colsStr := ""
		for _, el := range pairs {
			colsStr += ", `" + el[0].(string) + "` "
			values = append(values, el[1])
		}
		fmt.Println(values);
		queryStr = "INSERT INTO `" + name +
			"` (" + colsStr[1:] + ") VALUES (" + placeholders[1:] + ")"
	}
	return queryStr, values
}

func generateCreateTableString(pairs [][]interface{}) string {
	queryStr := ""
	colsStr := ""
	fkStr := ""
	tableName := ""
	for _, pair := range pairs {
		if pair[0] == "Table name" {
			tableName = pair[1].(string)
			queryStr = "CREATE TABLE `" + tableName +
				"` ( id INT NOT NULL AUTO_INCREMENT, "
		} else {
			col := pair[0].(string)
			val := pair[1].(string)
			if strings.HasSuffix(val, "_id") {
				colsStr += "`" + col + "_fk` INT, "
				fk := val[:len(val) - 3]
				fkStr += `,` + makeForeignKey(tableName, col, fk)
			} else {
				colsStr += " `" + col + "` " + val + ", "
			}
		}
	}
	queryStr += colsStr
	queryStr += ` PRIMARY KEY (id),
				UNIQUE INDEX id_UNIQUE (id ASC)
			` + fkStr + `
			)`
	return queryStr
}

func makeForeignKey(this, name, ref string) string {
	return `KEY fk_`+ this +`_`+ name + `_id_idx (`+ name +`_fk),
			  CONSTRAINT fk_`+ this +`_`+ name + `_id
			  FOREIGN KEY (`+ name +`_fk)
			  REFERENCES `+ ref +` (id) ON DELETE NO ACTION ON UPDATE NO ACTION`
}

func makeColsInfo(tableName, colName string) string {
	if colName != "" {
		return `select kcu.referenced_table_name
				from information_schema.key_column_usage kcu
				WHERE kcu.TABLE_SCHEMA   = 'superhero'
					AND  kcu.column_name   = '`+ colName +`'
					AND  kcu.table_name = '`+ tableName +`'`
	} else {
		return `SELECT cols.column_name, cols.column_type ,
				(select kcu.referenced_table_name from information_schema.key_column_usage kcu
		WHERE kcu.TABLE_SCHEMA   = 'superhero'
			AND  kcu.column_name   = cols.column_name
			AND  kcu.table_name = '`+ tableName +`' limit 1) as referenced_table_name
		FROM information_schema.columns cols
			WHERE cols.TABLE_SCHEMA = 'superhero' and cols.table_name = '`+ tableName +`';`
	}
}

// delete entity
func Delete(name string, id string) (string, error) {
	var err error
	if name == "home" {
		queryStr := "DROP TABLE " + id
		fmt.Println(queryStr)
		stmt, err := Db.Prepare(queryStr)
		CheckError(err)
		_, err = stmt.Exec()
		stmt.Close()
	} else {
		queryStr := "DELETE FROM " + name + " WHERE id = ?"
		fmt.Println(queryStr)
		stmt, err := Db.Prepare(queryStr)
		CheckError(err)
		_, err = stmt.Exec(id)
		stmt.Close()
	}
	CheckError(err)
	return "", err
}

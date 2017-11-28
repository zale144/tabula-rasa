package dao

import (
	. "tabula-rasa/util"
	"strings"
	"encoding/json"
	"fmt"
	"strconv"
)

// get all entities
func Get(args ...string) (string, error) {
	name := args[0]
	id := args[1]
	spec := args[2]
	single := false
	query := "SELECT * FROM " + name
	if id != "" && id != "0" {
		id = args[1]
		query += " WHERE id = " + id
		single = true
	} else if name == "home" {
		query = `SELECT DISTINCT table_name FROM information_schema.columns
		WHERE column_name in ('id') AND table_schema = 'superhero';`
	} else if spec == "cols" {
		query = `SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = 'superhero' AND TABLE_NAME = '` + name + "'"
	}

	rows, err := Db.Query(query)
	CheckError(err)
	cols, err := rows.Columns()
	CheckError(err)
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	var results string
	var objs []string

	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	for rows.Next() {
		err = rows.Scan(dest...)
		CheckError(err)
		for i, raw := range rawResult {
			result[i] = `"` + cols[i] + `" : "` + string(raw) + `"`
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

// create entity
func Create(name string, obj []byte) (string, error) {
	objMap := make(map[string]string)
	json.Unmarshal(obj, &objMap)
	cols := []string{}
	vals := []interface{}{}
	placeholders := []string{}
	var queryStr string
	var id int64 = 0
	for k, v := range objMap {
		cols = append(cols, k)
		vals = append(vals, v)
		if k == "Id" {
			id, _ = strconv.ParseInt(v, 0, 64)
		}
		placeholders = append(placeholders, "?")
	}
	if name != "home" {
		if id > 0 {
			queryStr = `UPDATE ` + name + ` SET `
			for i, el := range cols {
				if el != "Id" {
					queryStr += el + ` = '` + vals[i].(string) + `',`
				}
			}
			queryStr = queryStr[:len(queryStr) - 1] + ` WHERE id = ` + strconv.FormatInt(id, 10)
			vals = vals[:0]
		} else {
			colsStr := strings.Join(cols, ", ")
			plcStr := strings.Join(placeholders, ",")
			queryStr = "INSERT INTO " + name +
				" (" + colsStr + ") VALUES (" + plcStr + ")"
		}
	} else {
		colsStr := "";
		for i, val := range vals {
			if cols[i] != "Table name" {
				colsStr +=  val.(string) + ` VARCHAR(45),`
			} else {
				queryStr = `CREATE TABLE ` + vals[i].(string) +
					`(
						id INT NOT NULL AUTO_INCREMENT, `
			}
		}
		queryStr += colsStr
		queryStr += `PRIMARY KEY (id),
				UNIQUE INDEX id_UNIQUE (id ASC)
			)`
		vals = vals[:0]
	}
	stmt, err := Db.Prepare(queryStr)
	fmt.Println(queryStr)
	CheckError(err)
	defer stmt.Close()
	res, err := stmt.Exec(vals...)
	CheckError(err)
	id, _ = res.LastInsertId()
	return Get(name, fmt.Sprint(id), "")
}

// delete entity
func Delete(name string, id string) (string, error) {
	var err error
	if name == "home" {
		stmt, err := Db.Prepare("DROP TABLE " + id)
		CheckError(err)
		_, err = stmt.Exec()
		stmt.Close()
	} else {
		stmt, err := Db.Prepare("DELETE FROM " + name + " WHERE id = ?")
		CheckError(err)
		_, err = stmt.Exec(id)
		stmt.Close()
	}
	CheckError(err)
	return "", err
}

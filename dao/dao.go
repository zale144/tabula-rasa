package dao

import (
	. "simple_server/util"
	"strings"
	"encoding/json"
	"fmt"
)

// get all entities
func Get(args ...string) (string, error) {
	name := args[0]
	id := args[1]
	spec := args[2]
	single := false
	query := "SELECT * FROM " + name
	if id != "" {
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
	for k, v := range objMap {
		cols = append(cols, k)
		vals = append(vals, v)
		placeholders = append(placeholders, "?")
	}
	colsStr := strings.Join(cols, ", ")
	plcStr := strings.Join(placeholders, ",")
	queryStr := "INSERT INTO " + name +
		" (" + colsStr + ") VALUES (" + plcStr + ")"
	stmt, err := Db.Prepare(queryStr)
	CheckError(err)
	defer stmt.Close()
	res, err := stmt.Exec(vals...)
	CheckError(err)
	id, _ := res.LastInsertId()
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

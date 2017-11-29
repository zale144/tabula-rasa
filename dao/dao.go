package dao

import (
	. "tabula-rasa/util"
	"strings"
	"fmt"
	"strconv"
)

func mapParameters(args []string) (string, string, string, string) {
	name, id, spec, col := "", "", "", ""
	params := []string{name, id, spec, col}
	for i, arg := range args {
		params[i] = arg
	}
	return params[0], params[1], params[2], params[3]
}

// get all entities
func Get(args ...string) (string, error) {
	name, id, spec, col := mapParameters(args)
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
		query = makeColsInfo(name, col)
	}
	fmt.Println(query)
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

// create entity
func Create(name string, obj []byte) (string, error) {
	pairs, id := unmarshalJson(obj)
	var queryStr string
	var values []interface{}

	if name != "home" {
		queryStr, values = generateInsertUpdateString(name , id, pairs)
	} else {
		queryStr = generateCreateTableString(pairs)
	}
	stmt, err := Db.Prepare(queryStr)
	fmt.Println(queryStr)
	CheckError(err)
	defer stmt.Close()
	res, err := stmt.Exec(values...)
	CheckError(err)
	id, _ = res.LastInsertId()
	return Get(name, fmt.Sprint(id), "")
}

func unmarshalJson(obj []byte) ([][]interface{}, int64) {
	var id int64
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
		if pair[0] == "Id" {
			id, _ = strconv.ParseInt(pair[1], 0, 64)
		}
	}
	return pairList, id
}

func generateInsertUpdateString(name string, id int64, pairs [][]interface{}) (string, []interface{}) {
	placeholders := generatePlaceholders(len(pairs))
	values := []interface{}{}
	queryStr := ""
	if id > 0 {
		queryStr = `UPDATE ` + name + ` SET `
		for _, el := range pairs {
			if el[0] != "Id" {
				queryStr += el[0].(string) + ` = '` + el[1].(string) + `',`
			}
		}
		queryStr = queryStr[:len(queryStr) - 1] + ` WHERE id = ` + strconv.FormatInt(id, 10)
		pairs = pairs[:0]
	} else {
		colsStr := ""
		for _, el := range pairs {
			if el[0] != "Id" {
				colsStr += ", `" + el[0].(string) + "` "
			}
			values = append(values, el[1])
		}
		queryStr = "INSERT INTO " + name +
			" (" + colsStr[1:] + ") VALUES (" + placeholders + ")"
	}
	return queryStr, values;
}

func generateCreateTableString(pairs [][]interface{}) string {
	queryStr := ""
	colsStr := ""
	fkStr := ""
	fk := ""
	tableName := ""
	for _, pair := range pairs {
		if pair[0] == "Table name" {
			tableName = pair[1].(string)
			queryStr = `CREATE TABLE ` + tableName +
				`(	id INT NOT NULL AUTO_INCREMENT, `
		} else {
			col := pair[0].(string)
			val := pair[1].(string)
			if strings.HasSuffix(val, "_id") {
				colsStr += col + `_fk INT, `
				fk = val[:len(val) - 3]
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

func generatePlaceholders(num int) string {
	placehoders := ""
	for i := 0; i < num; i++ {
		placehoders += ",?"
	}
	return placehoders[1:]
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

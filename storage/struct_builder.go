package storage

import (
	"reflect"
	"strings"
	"strconv"
	"log"
)

type StructBuilder struct {}

// convert a map representation of a row in the database to
// a new custom generated struct with appropriate
// field names and types
func (sb StructBuilder) convertRowMapToStruct(cols []ColumnInfo, typeNames []string, m map[string]interface{}) interface{} {
	var structFields []reflect.StructField
	// iterate through all the column names
	for i := range cols {
		// set the name and type of the struct field
		field := reflect.StructField{
			Name: strings.Title(cols[i].Name), // capitalize, so json can marshal it
			Type: reflect.TypeOf(typeMap[typeNames[i]]),
		}
		// append the field to the struct
		structFields = append(structFields, field)
	}
	structType := reflect.StructOf(structFields)
	// get a reflected pointer to the new struct type
	structValue := reflect.New(structType)
	// create the type that will hold the row data
	obj := structValue.Interface()
	sr := reflect.ValueOf(obj)
	// iterate through the column names
	for i, c := range cols {
		if m[c.Name] == nil {
			continue
		}
				// set value of appropriate type to struct field
		switch typeNames[i] {
		// for string types
		case Varchar, Text:
			strVal := string(m[c.Name].([]uint8))
			sr.Elem().Field(i).SetString(strVal)
		// for integer types
		case Int, Integer:
			intVal, err := strconv.ParseInt(string(m[c.Name].([]uint8)), 0, 64)
			if err != nil {
				log.Fatal(err)
				return err
			}
			sr.Elem().Field(i).SetInt(intVal)
			/*if c.Reference != "" {
				fmt.Printf( "Column %v with 'id' %v has a reference to table %v\n", c, intVal, c.Reference)
			}*/
		// for floating point types
		case Double, Float:
			dblVal, err := strconv.ParseFloat(string(m[c.Name].([]uint8)),64)
			if err != nil {
				log.Fatal(err)
				return err
			}
			sr.Elem().Field(i).SetFloat(dblVal)
		// for boolean types
		case Boolean, Bool:
			boolVal, err := strconv.ParseBool(string(m[c.Name].([]uint8)))
			if err != nil {
				log.Fatal(err)
				return err
			}
			sr.Elem().Field(i).SetBool(boolVal)
		}
	}
	return sr.Interface()
}

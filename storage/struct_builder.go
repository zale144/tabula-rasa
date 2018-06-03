package storage

import (
	"reflect"
	"strings"
	"log"
)

type StructBuilder struct {
	Struct reflect.Value
}

/*-----------------------------------------------------------
 * convert a map representation of a row in the database to *
 * a new custom generated struct with appropriate		    *
 * field names and types, and fill the data				    *
 *----------------------------------------------------------*/
func (sb StructBuilder) convertRowMapToStruct(cols []ColumnInfo, m map[string]interface{}) interface{} {
	// create a custom struct that will hold the data from the row map
	sb.createCustomStruct(cols)
	// iterate through the column names
	for i, c := range cols {
		if m[*c.Name] == nil {
			continue
		}
		// use the Set method to set the field value, according to type
		err := TypeSetter{}.Set(m[*c.Name], sb.Struct.Elem().Field(i), *cols[i].Type)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return sb.Struct.Interface()
}
// method for creating a custom struct from predefined column names and types
func (sb *StructBuilder) createCustomStruct(cols []ColumnInfo)  {
	var structFields []reflect.StructField
	// iterate through all the column names
	for i := range cols {
		// if column is a struct itself, get it's type and
		// set the new custom struct's type to it
		var val interface{}
		if cols[i].Child != nil {
			val = *cols[i].Child
		}
		// use the type setter to get the type for the new struct field
		typ := TypeSetter{}.Set(val, nil, *cols[i].Type)
		// set the name and type of the struct field
		field := reflect.StructField{
			Name: strings.Title(*cols[i].Name), // capitalize, so json can marshal it
			Type: reflect.TypeOf(typ),
		}
		// append the field to the struct
		structFields = append(structFields, field)
	}
	structType := reflect.StructOf(structFields)
	// get a reflected pointer to the new struct type
	structValue := reflect.New(structType)
	// create the type that will hold the row data
	obj := structValue.Interface()
	sb.Struct = reflect.ValueOf(obj)
}

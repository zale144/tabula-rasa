package storage

import (
	"reflect"
	"log"
	"strconv"
)

// declarations of the database data types
const (
	Text = "TEXT"
	Varchar = "VARCHAR"
	Int = "INT"
	Integer = "INTEGER"
	Double = "DOUBLE"
	Float = "FLOAT"
	Boolean = "BOOLEAN"
	Bool = "BOOL"
	SingleRef = "SINGLE_REFERENCE"
)
// a map holding all of the data type reference functions that we will use
// while setting field type/value during custom struct creation
var (
	typeSetter = map[string]interface{}{
		Text: setString,
		Varchar: setString,
		Int: setInteger,
		Integer: setInteger,
		Double: setDouble,
		Float: setDouble,
		Boolean: setBoolean,
		Bool: setBoolean,
		SingleRef: setSingleReference,
	}
)

type TypeSetter struct {
	SetterFunc func(val interface{}, field interface{})interface{}
}

// set value/type according to value, field, type
func (ts TypeSetter) Set(val interface{}, field interface{}, typ string) interface{} {
	// set the setter func by type name
	ts.setFunc(typ)
	// return the error/type
	return ts.SetterFunc(val, field)
}
// method for setting the setter function according to data type name
func (ts *TypeSetter) setFunc(typ string) {
	// get the setter func from the typeSetter map
	ts.SetterFunc = typeSetter[typ].(func(val interface{}, field interface{})interface{})
}
// function for setting a string value/type to a struct field
func setString(val interface{}, field interface{}) interface{} {
	if field == nil {
		return *new(string)
	}
	strVal := string(val.([]uint8))
	field.(reflect.Value).SetString(strVal)
	return nil
}
// function for setting an integer value/type to a struct field
func setInteger(val interface{}, field interface{}) interface{} {
	if field == nil {
		return *new(int)
	}
	strVal := string(val.([]uint8))
	intVal, err := strconv.ParseInt(strVal, 0, 64)
	if err != nil {
		log.Fatal(err)
		return err
	}
	field.(reflect.Value).SetInt(intVal)
	return nil
}
// function for setting a double value/type to a struct field
func setDouble(val interface{}, field interface{}) interface{} {
	if field == nil {
		return *new(float64)
	}
	strVal := string(val.([]uint8))
	dblVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Fatal(err)
		return err
	}
	field.(reflect.Value).SetFloat(dblVal)
	return nil
}
// function for setting a double value/type to a struct field
func setBoolean(val interface{}, field interface{}) interface{} {
	if field == nil {
		return *new(bool)
	}
	strVal := string(val.([]uint8))
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		log.Fatal(err)
		return err
	}
	field.(reflect.Value).SetBool(boolVal)
	return nil
}
// function for setting a single reference value/type to a struct field
func setSingleReference(val interface{}, field interface{}) interface{} {
	if field == nil {
		return reflect.ValueOf(val).Interface()
	}
	field.(reflect.Value).Set(reflect.ValueOf(val))
	return nil
}
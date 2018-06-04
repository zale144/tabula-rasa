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

type TypeSetter struct {
	value interface{}
}

type Type interface {
	Set(val interface{}, field interface{}) (interface{}, error)
}
// a map holding all of the data type reference functions that we will use
// while setting field type/value during custom struct creation
var (
	typeMap = map[string]Type{
		Text:    new(text),
		Varchar: new(text),
		Int: new(integer),
		Integer: new(integer),
		Double: new(double),
		Float: new(double),
		Boolean: new(boolean),
		Bool: new(boolean),
		SingleRef: new(singleReference),
	}
)
// set value/type according to value, field, type
func (ts *TypeSetter) SetType(val interface{}, field interface{}, typ string) error {
	typeStruct := typeMap[typ]
	v, err := typeStruct.Set(val, field)
	ts.value = v
	return err
}
// implements Type
type text struct {}
// method for setting a string value/type to a struct field
func (s text) Set(val interface{}, field interface{}) (interface{}, error) {
	if field == nil {
		return *new(string), nil
	}
	strVal := string(val.([]uint8))
	field.(reflect.Value).SetString(strVal)
	return nil, nil
}
// implements Type
type integer struct {}
// method for setting an integer value/type to a struct field
func (i integer) Set(val interface{}, field interface{}) (interface{}, error) {
	if field == nil {
		return *new(int), nil
	}
	strVal := string(val.([]uint8))
	intVal, err := strconv.ParseInt(strVal, 0, 64)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	field.(reflect.Value).SetInt(intVal)
	return nil, nil
}
// implements Type
type double struct {}
// method for setting a double value/type to a struct field
func (d double) Set(val interface{}, field interface{}) (interface{}, error) {
	if field == nil {
		return *new(float64), nil
	}
	strVal := string(val.([]uint8))
	dblVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	field.(reflect.Value).SetFloat(dblVal)
	return nil, nil
}
// implements Type
type boolean struct {}
// method for setting a double value/type to a struct field
func (b boolean) Set(val interface{}, field interface{}) (interface{}, error) {
	if field == nil {
		return *new(bool), nil
	}
	strVal := string(val.([]uint8))
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	field.(reflect.Value).SetBool(boolVal)
	return nil, nil
}
// implements Type
type singleReference struct {}
// method for setting a single reference value/type to a struct field
func (sr singleReference) Set(val interface{}, field interface{}) (interface{}, error) {
	if field == nil {
		return interface{}(reflect.ValueOf(val).Interface()), nil
	}
	field.(reflect.Value).Set(reflect.ValueOf(val))
	return nil, nil
}
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
	typeMap = map[string]Type{
		Text:    *new(text),
		Varchar: *new(text),
		Int: *new(integer),
		Integer: *new(integer),
		Double: *new(double),
		Float: *new(double),
		Boolean: *new(boolean),
		Bool: *new(boolean),
		SingleRef: new(singleReference),
	}
)

// get value of type according to value, type
func GetType(val interface{}, typ string) interface{} {
	typeStruct := typeMap[typ]
	// if type is a child reference
	if val != nil {
		sr := typeStruct.(*singleReference)
		sr.value = val
	}
	return typeStruct.Get()
}
// set type value according to value, field, type
func SetType(val interface{}, field interface{}, typ string)  {
	typeMap[typ].Set(val, field)
}

type Type interface {
	Get()interface{}
	Set(val interface{}, field interface{})
}

/*************************
* String type	 		 *
* - implements Type		 *
*************************/
type text string

// method for getting a string type value
func (s text) Get() interface{} {
	return string(s)
}

// method for setting a string type value to a struct field
func (s text) Set(val interface{}, field interface{})  {
	strVal := string(val.([]uint8))
	field.(reflect.Value).SetString(strVal)
}

/*************************
* Integer type		 	 *
* - implements Type		 *
*************************/
type integer int

// method for getting an integer type value
func (i integer) Get() interface{} {
	return int(i)
}

// method for setting an integer value/type to a struct field
func (i integer) Set(val interface{}, field interface{})  {
	strVal := string(val.([]uint8))
	intVal, err := strconv.ParseInt(strVal, 0, 64)
	if err != nil {
		log.Fatal(err)
	}
	field.(reflect.Value).SetInt(intVal)
}

/*************************
* Floating point type	 *
* - implements Type		 *
*************************/
type double float64

// method for getting a floating point type value
func (d double) Get() interface{} {
	return float64(d)
}

// method for setting a floating point type value to a struct field
func (d double) Set(val interface{}, field interface{})  {
	strVal := string(val.([]uint8))
	dblVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Fatal(err)
	}
	field.(reflect.Value).SetFloat(dblVal)
}

/*************************
* Boolean type	 		 *
* - implements Type		 *
*************************/
type boolean bool

// method for getting a boolean type value
func (b boolean) Get() interface{} {
	return bool(b)
}

// method for setting a boolean type value to a struct field
func (b boolean) Set(val interface{}, field interface{})  {
	strVal := string(val.([]uint8))
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		log.Fatal(err)
	}
	field.(reflect.Value).SetBool(boolVal)
}

/*************************
* Single Reference type	 *
* - implements Type		 *
*************************/
type singleReference struct {
	value interface{}
}

// method for getting a single reference type value
func (sr singleReference) Get() interface{} {
	return interface{}(reflect.ValueOf(sr.value).Interface())
}

// method for setting a single reference type value to a struct field
func (sr singleReference) Set(val interface{}, field interface{})  {
	field.(reflect.Value).Set(reflect.ValueOf(val))
}
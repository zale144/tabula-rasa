package storage

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
)
// a map holding all of the data type references that we will use
var (
	typeMap = map[string]interface{}{
		Text: *new(string),
		Varchar: *new(string),
		Int: *new(int),
		Integer: *new(int),
		Double: *new(float64),
		Float: *new(float64),
		Boolean: *new(bool),
		Bool: *new(bool),
	}
)
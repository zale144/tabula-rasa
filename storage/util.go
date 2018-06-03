package storage

import (
	"reflect"
	"log"
)
// method for getting complete column info
func (ts TableStorage) getColumnInfo(name, typ, id string, query *string) ([]ColumnInfo, error) {
	var columnInfo = []ColumnInfo{}
	if name != "tables" && typ == "rows" {
		// call the Get method to get the column information for this table
		columns, err := ts.Get(name, "", "cols")
		if err != nil {
			return nil, err
		}
		columnInfo = make([]ColumnInfo, len(columns))
		// retrieve the reference info for the column if any
		err = ts.getReferenceInfo(&columnInfo, columns, query, id, typ, name)
		if err != nil {
			return nil, err
		}
	}
	return columnInfo, nil
}
// method for retrieving referenced tables info
func (ts TableStorage) getReferenceInfo(columnInfo *[]ColumnInfo, columns []interface{}, query *string, id, typ, name string) error {
	// get the table column info
	colInfo := *columnInfo
	for i := range columns {
		col := columns[i]
		elem := reflect.ValueOf(col).Elem()
		fld := elem.Field(0) // name of the column
		ref := elem.Field(2) // name of the referenced table
		colName := fld.Interface()
		reference := ref.Interface()
		colNameStr := colName.(string)
		colInfo[i].Name = &colNameStr
		refStr := reference.(string)
		colInfo[i].Reference = &refStr
		if *colInfo[i].Reference != "" {
			*query = generateSelectQuery(name, id, typ)
		}
	}
	return nil
}
// method for adding values to the row map
func (ts TableStorage) addValuesToRowMap(vals *[]interface{}, m *map[string]interface{}, columnInfo *[]ColumnInfo) {
	//refMapV := *refMap
	mV := *m
	colInfo := *columnInfo
	for i, val := range *vals {
		mV[*colInfo[i].Name] = val
		// if there is a referenced table
		if colInfo[i].Reference != nil && val != nil {
			ref := *colInfo[i].Reference
			if ref == "" {
				continue
			}
			// get the referenced struct by it's id
			reference, err := ts.Get(ref, string(val.([]uint8)), "rows")
			if err != nil {
				log.Println(err)
				continue
			}
			refStruct := reflect.ValueOf(reference[0]).Elem().Interface()
			mV[*colInfo[i].Name] = refStruct
			colInfo[i].Child = &refStruct
			// TODO many-to-many reference
			sRef := SingleRef // set the type of reference
			colInfo[i].Type = &sRef
		}
	}
}
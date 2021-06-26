// @Author: abbeymart | Abi Akindele | @Created: 2020-12-09 | @Updated: 2020-12-09
// @Company: mConnect.biz | @License: MIT
// @Description: crud utility functions

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"github.com/asaskevich/govalidator"
	"reflect"
)

type EmailUserNameType struct {
	Email    string
	Username string
}

func EmailUsername(loginName string) EmailUserNameType {
	if govalidator.IsEmail(loginName) {
		return EmailUserNameType{
			Email:    loginName,
			Username: "",
		}
	}

	return EmailUserNameType{
		Email:    "",
		Username: loginName,
	}

}

func ParseRawValues(rawValues [][]byte) ([]interface{}, error) {
	// variables
	var v interface{}
	var va []interface{}
	// parse the current-raw-values
	for _, val := range rawValues {
		if err := json.Unmarshal(val, &v); err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing raw-row-value: %v", err.Error()))
		} else {
			va = append(va, v)
		}
	}
	return va, nil
}

func ArrayStringContains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ArrayIntContains(arr []int, val int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ArraySQLInStringValues(arr []string) string {
	result := ""
	for ind, val := range arr {
		result += "'" + val + "'"
		if ind < len(arr)-1 {
			result += ", "
		}
	}
	return result
}

// JsonDataETL method converts json inputs to equivalent struct data type specification
// rec must be a pointer to a type matching the jsonRec
func JsonDataETL(jsonRec []byte, rec interface{}) error {
	if err := json.Unmarshal(jsonRec, &rec); err == nil {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error converting json-to-record-format: %v", err.Error()))
	}
}

// DataToValueParam method accepts only a struct record/param (type/model) and returns the ActionParamType
// data camel/Pascal-case keys are converted to underscore-keys to match table-field/columns specs
func DataToValueParam(rec interface{}) (types.ActionParamType, error) {
	dataValue := types.ActionParamType{}
	v := reflect.ValueOf(rec)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		dataValue[govalidator.CamelCaseToUnderscore(typeOfS.Field(i).Name)] = v.Field(i).Interface()
		//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	return dataValue, nil
}

func DataToValueParam2(rec interface{}) (types.ActionParamType, error) {
	switch rec.(type) {
	case struct{}:
		dataValue := types.ActionParamType{}
		v := reflect.ValueOf(rec)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			dataValue[govalidator.CamelCaseToUnderscore(typeOfS.Field(i).Name)] = v.Field(i).Interface()
			//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		}
		return dataValue, nil
	default:
		return nil, errors.New("invalid type - requires parameter of type struct only")
	}
}

// StructToMap function converts struct to map
func StructToMap(rec interface{}) map[string]interface{} {
	var mapDataValue map[string]interface{}
	jsonRec, _ := json.Marshal(rec)
	err := json.Unmarshal(jsonRec, &mapDataValue)
	if err != nil {
		return nil
	}
	return mapDataValue
}

// TagField return the field-tag (e.g. table-column-name) for mcorm tag
func TagField(rec interface{}, fieldName string, tag string) string {
	t := reflect.TypeOf(rec)
	field, found := t.FieldByName(fieldName)
	if !found {
		return ""
	}
	//tagValue := field.Tag
	return field.Tag.Get(tag)
}

// StructToTagMap function converts struct to map
func StructToTagMap(rec interface{}, tag string) map[string]interface{} {
	tagMapDataValue := map[string]interface{}{}
	mapDataValue := StructToMap(rec)
	// compose tagMapDataValue
	for key, val := range mapDataValue {
		tagField := TagField(rec, key, tag)
		tagMapDataValue[tagField] = val
	}
	return tagMapDataValue
}

// StructToFieldValues function converts struct to map
func StructToFieldValues(rec interface{}, tag string) ([]string, []interface{}, error) {
	var tableFields []string
	var fieldValues []interface{}
	// validate rec as a struct{}
	switch rec.(type) {
	case struct{}:
		mapDataValue := StructToMap(rec.(struct{}))
		// compose tagMapDataValue
		for key, val := range mapDataValue {
			tagField := TagField(rec.(struct{}), key, tag)
			tableFields = append(tableFields, tagField)
			fieldValues = append(fieldValues, val)
		}
		return tableFields, fieldValues, nil
	default:
		return nil, nil, errors.New("invalid type - requires parameter of type struct only")
	}

}

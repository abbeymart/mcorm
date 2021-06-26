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
	"strings"
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
func StructToMap(rec interface{}) (map[string]interface{}, error) {
	var mapData map[string]interface{}
	// json record
	jsonRec, err := json.Marshal(rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	// json-to-map
	err = json.Unmarshal(jsonRec, &mapData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	return mapData, nil
}

// TagField return the field-tag (e.g. table-column-name) for mcorm tag
func TagField(rec interface{}, fieldName string, tag string) (string, error) {
	// TODO: validate rec as struct{}
	t := reflect.TypeOf(rec)
	// convert the first-letter to upper-case (public field)
	field, found := t.FieldByName(strings.Title(fieldName))
	if !found {
		// check private field
		field, found = t.FieldByName(fieldName)
		if !found {
			return "", errors.New(fmt.Sprintf("error retrieving tag-field for field-name: %v", fieldName))
		}
	}
	//tagValue := field.Tag
	return field.Tag.Get(tag), nil
}

// StructToTagMap function converts struct to map (for crud-actionParams / records)
func StructToTagMap(rec interface{}, tag string) (map[string]interface{}, error) {
	tagMapData := map[string]interface{}{}
	mapData, err := StructToMap(rec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error computing struct to map: %v", err.Error()))
	}
	// compose tagMapData
	for key, val := range mapData {
		tagField, tagErr := TagField(rec, key, tag)
		if tagErr != nil {
			return nil, errors.New(fmt.Sprintf("error computing tag-field: %v", tagErr.Error()))
		}
		tagMapData[tagField] = val
	}
	return tagMapData, nil
}

// StructToFieldValues function converts struct/map to map (for DB columns and values)
func StructToFieldValues(rec interface{}, tag string) ([]string, []interface{}, error) {
	var tableFields []string
	var fieldValues []interface{}
	mapDataValue, err := StructToMap(rec)
	if err != nil {
		return nil, nil, errors.New("error computing struct to map")
	}
	// compose tagMapDataValue
	for key, val := range mapDataValue {
		tagField, tagErr := TagField(rec.(struct{}), key, tag)
		if tagErr != nil {
			return nil, nil, errors.New(fmt.Sprintf("error retrieving tag-field: %v", key))
		}
		tableFields = append(tableFields, tagField)
		fieldValues = append(fieldValues, val)
	}
	return tableFields, fieldValues, nil
}

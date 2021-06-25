// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"github.com/asaskevich/govalidator"
	"reflect"
)

// data
var jsonData = `[
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"},
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"},
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"}
]`

var jsonQueryParams = `
	{
		"group_name": "abc",
		"group_order": 1,
		"group_link_op": "and",
		"group_items": [
			{"group_item": {"name": { "eq": Paul"}},
			"group_item_order": 1,
			"groupItemOp": "and"
			},
			{"group_item": {"age":{"gt": 10}},
			"group_item_order": 3,
			},
			{"group_item": {"location": {"eq": Toronto"}}
			"group_item_order": 2,
			"groupItemOp": "and"
			},
		]
	},
	{},
	{},
`

// convert/decode jsonQueryParams to queryParams
var queryParams types.QueryParamType
var _ = json.Unmarshal([]byte(jsonQueryParams), &queryParams)
//var _ = JsonDataETL([]byte(jsonQueryParams), queryParams)

type Person struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	LocationId  string `json:"location_id"`
	PhoneNumber string `json:"phone_number"`
}

//var camelCaseUnderscore = govalidator.CamelCaseToUnderscore("firstName")

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

// convert dataToValueParam record(s) to action-params
var actionParams = types.ActionParamsType{
	types.ActionParamType{},
	types.ActionParamType{},
	types.ActionParamType{},
	types.ActionParamType{},
}

func main() {
	var person Person
	if err := JsonDataETL([]byte(jsonData), &person); err == nil {
		fmt.Printf("Person's record: %+v", person)
	} else {
		fmt.Printf("Error coverting json-data: %v", err.Error())
	}

	// TODO: table-fields | tableFieldPointers | queryParams

}

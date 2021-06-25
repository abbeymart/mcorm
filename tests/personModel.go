// @Author: abbeymart | Abi Akindele | @Created: 2020-12-15 | @Updated: 2020-12-15
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import (
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mctypes/datatypes"
)

var PersonModel = types.ModelType{
	TableName: "persons",
	RecordDesc: map[string]types.FieldDescType{
		"id": {
			FieldType: datatypes.String,
			FieldLength: 100,
			FieldPattern: "",
			AllowNull: false,
			Unique: false,
			Indexable: false,
			PrimaryKey: false,
			ValidateMessage: "Length must not be longer than 100",
		},
		"name": {
			FieldType: datatypes.String,
			FieldLength: 100,
			FieldPattern: "",
			AllowNull: false,
			Unique: false,
			Indexable: false,
			PrimaryKey: false,
			ValidateMessage: "Length must not be longer than 100",
		},
	},
	Relations: nil,
	TimeStamp: true,
	ActorStamp: true,
	ActiveStamp: true,
	AlterSyncTable: true,
}

//var person = mccrud.NewModel(PersonModel)

// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: get/read records test cases

package mcorm

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/datatypes"
	"github.com/abbeymart/mctest"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	myDb := mcdb.DbConfig{
		DbType:   "postgres",
		Host:     "localhost",
		Username: "postgres",
		Password: "ab12testing",
		Port:     5432,
		DbName:   "mcdev",
		Filename: "testdb.db",
		PoolSize: 20,
		Url:      "localhost:5432",
	}
	myDb.Options = mcdb.DbConnectOptions{}

	type AuditType struct {
		Id            string      `json:"id"`
		TableName     string      `json:"tableName"`
		LogRecords    interface{} `json:"logRecords"`
		NewLogRecords interface{} `json:"newLogRecords"`
		LogType       string      `json:"logType"`
		LogBy         string      `json:"logBy"`
		LogAt         time.Time   `json:"logAt"`
		//dbcrud.AuditStampType
	}

	// db-connection
	dbc, err := myDb.OpenPgxDbPool()
	// defer dbClose
	defer myDb.ClosePgxDbPool()

	// check db-connection-error
	if err != nil {
		fmt.Printf("*****db-connection-error: %v\n", err.Error())
		return
	}

	auditModel := types.ModelType{
		TableName: "persons",
		RecordDesc: map[string]types.FieldDescType{
			"id": {
				FieldType:       datatypes.String,
				FieldLength:     100,
				FieldPattern:    "",
				AllowNull:       false,
				Unique:          false,
				Indexable:       false,
				PrimaryKey:      false,
				ValidateMessage: "Length must not be longer than 100",
			},
			"name": {
				FieldType:       datatypes.String,
				FieldLength:     100,
				FieldPattern:    "",
				AllowNull:       false,
				Unique:          false,
				Indexable:       false,
				PrimaryKey:      false,
				ValidateMessage: "Length must not be longer than 100",
			},
		},
		Relations:      nil,
		TimeStamp:      true,
		ActorStamp:     true,
		ActiveStamp:    true,
		AlterSyncTable: true,
	}

	AuditModel := NewModel(auditModel)

	getCrudParams := types.CrudParamsType{
		AppDb:       dbc.DbConn,
		TableName:   TestTable,
		UserInfo:    TestUserInfo,
		RecordIds:   GetIds,
		QueryParams: GetParams,
	}

	var getCrud = NewCrud(getCrudParams, TestCrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Ids and return success:",
		TestFunc: func() {
			logRec := AuditType{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			fmt.Printf("get-by-id-response: %#v\n\n", res)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-id-value: %#v\n", value.TableRecords)
			fmt.Printf("get-by-param-count: %v\n", value.RecordCount)
			jsonRecs, _ := json.Marshal(value.TableRecords)
			fmt.Printf("json-records: %v\n\n", string(jsonRecs))
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "get-task-count should be: 2")
			mctest.AssertEquals(t, len(value.TableRecords), 2, "get-result-count should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get records by query-params and return success:",
		TestFunc: func() {
			logRec := AuditType{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			//fmt.Printf("get-by-param-response: %#v\n", res)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-param-value: %#v\n", value.TableRecords)
			fmt.Printf("get-by-param-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount >= 0, true, "get-task-count should be >= 0")
			mctest.AssertEquals(t, len(value.TableRecords) >= 0, true, "get-result-count should be >= 0")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get all records and return success:",
		TestFunc: func() {
			logRec := AuditType{}
			getCrudParams.QueryParams = types.QueryParamType{}
			getCrudParams.RecordIds = []string{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount >= 10, true, "get-task-count should be >= 10")
			mctest.AssertEquals(t, len(value.TableRecords) >= 10, true, "get-result-count should be >= 10")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records by limit/skip(offset) and return success:",
		TestFunc: func() {
			logRec := AuditType{}
			getCrud.Skip = 0
			getCrud.Limit = 20
			getCrudParams.QueryParams = types.QueryParamType{}
			getCrudParams.RecordIds = []string{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-limit-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount == 20, true, "get-task-count should be = 20")
			mctest.AssertEquals(t, len(value.TableRecords) == 20, true, "get-result-count should be = 20")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Id and return success[get-record method]:",
		TestFunc: func() {
			logRec := AuditType{}
			getCrudParams.RecordIds = GetIds
			getCrudParams.QueryParams = types.QueryParamType{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			//fmt.Printf("get-by-param-response: %#v\n", res)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-all-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "get-task-count should be 2")
			mctest.AssertEquals(t, len(value.TableRecords), 2, "get-result-count should be 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get records by params and return success[get-record method]:",
		TestFunc: func() {
			logRec := AuditType{}
			getCrudParams.RecordIds = []string{}
			getCrudParams.QueryParams = GetParams
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-limit-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount > 0, true, "get-task-count should be > 0")
			mctest.AssertEquals(t, len(value.TableRecords) > 0, true, "get-result-count should be > 0")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records and return success[get-record method]:",
		TestFunc: func() {
			logRec := AuditType{}
			getCrudParams.Skip = 0
			getCrudParams.Limit = 20
			getCrudParams.RecordIds = []string{}
			getCrudParams.QueryParams = types.QueryParamType{}
			res := AuditModel.Get(logRec, getCrudParams, TestCrudParamOptions)
			value, _ := res.Value.(types.CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-limit-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount == 20, true, "get-task-count should be = 20")
			mctest.AssertEquals(t, len(value.TableRecords) == 20, true, "get-result-count should be = 20")
		},
	})

	mctest.PostTestResult()

}

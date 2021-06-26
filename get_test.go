// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: get/read records test cases

package mcorm

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mcorm/types"
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
			var (
				//id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			var logRec AuditType
			//tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			tableFieldPointers := []interface{}{&logRec.Id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			res := getCrud.GetById(GetTableFields, tableFieldPointers)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			res := getCrud.GetByParam(GetTableFields, tableFieldPointers)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			res := getCrud.GetAll(GetTableFields, tableFieldPointers)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			getCrud.Skip = 0
			getCrud.Limit = 20
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			res := getCrud.GetAll(GetTableFields, tableFieldPointers)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			getCrud.RecordIds = GetIds
			getCrud.QueryParams = types.QueryParamType{}
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			getRecParams := types.GetCrudParamsType{
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := getCrud.GetRecord(getRecParams)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			getCrud.RecordIds = []string{}
			getCrud.QueryParams = GetParams
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			getRecParams := types.GetCrudParamsType{
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := getCrud.GetRecord(getRecParams)
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
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)

			getCrud.Skip = 0
			getCrud.Limit = 20
			getCrud.RecordIds = []string{}
			getCrud.QueryParams = types.QueryParamType{}
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			getRecParams := types.GetCrudParamsType{
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := getCrud.GetRecord(getRecParams)
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

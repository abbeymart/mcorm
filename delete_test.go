// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: records deletion test cases

package mcorm

import (
	"fmt"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mctest"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
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

	// db-connection
	dbc, err := myDb.OpenPgxDbPool()
	// defer dbClose
	defer myDb.ClosePgxDbPool()
	// check db-connection-error
	if err != nil {
		fmt.Printf("*****db-connection-error: %v\n", err.Error())
		return
	}
	deleteCrudParams := types.CrudParamsType{
		AppDb:       dbc.DbConn,
		TableName:   TestTable,
		UserInfo:    TestUserInfo,
		RecordIds:   DeleteIds,
		QueryParams: DeleteParams,
	}
	deleteAllCrudParams := types.CrudParamsType{
		AppDb:     dbc.DbConn,
		TableName: DeleteAllTable,
		UserInfo:  TestUserInfo,
	}

	var deleteCrud = NewCrud(deleteCrudParams, TestCrudParamOptions)
	var deleteAllCrud = NewCrud(deleteAllCrudParams, TestCrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by Ids and return success:",
		TestFunc: func() {
			res := deleteCrud.DeleteById()
			fmt.Printf("delete-by-ids: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-id should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by query-params and return success:",
		TestFunc: func() {
			res := deleteCrud.DeleteByParam()
			fmt.Printf("delete-by-params: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-params should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete all table records and return success:",
		TestFunc: func() {
			res := deleteAllCrud.DeleteAll()
			fmt.Printf("delete-all: %v : %v \n", res.Message, res.ResCode)
			value := res.Value
			deleted, _ := value.(bool)
			mctest.AssertEquals(t, res.Code, "success", "delete-all should return code: success")
			mctest.AssertEquals(t, deleted, true, "deleted() must be true")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by Ids, log-task, and return success:",
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
			res := deleteCrud.DeleteByIdLog(DeleteSelectTableFields, tableFieldPointers)
			fmt.Printf("delete-by-ids-log: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-id-log should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by query-params, log-task and return success:",
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
			res := deleteCrud.DeleteByParamLog(DeleteSelectTableFields, tableFieldPointers )
			fmt.Printf("delete-by-params-log: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-params-log should return code: success")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by Ids and return success[delete-record-method]:",
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
			deleteCrud.RecordIds = DeleteIds
			deleteCrud.QueryParams = types.QueryParamType{}
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			deleteRecParams := types.DeleteCrudParamsType{
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := deleteCrud.DeleteRecord(deleteRecParams)
			fmt.Printf("delete-by-ids[delete-record]: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-id should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by query-params and return success[delete-record-method]:",
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
			deleteCrud.RecordIds = []string{}
			deleteCrud.QueryParams = DeleteParams
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			deleteRecParams := types.DeleteCrudParamsType{
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := deleteCrud.DeleteRecord(deleteRecParams)
			fmt.Printf("delete-by-params[delete-record]: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "delete-by-params-log should return code: success")
		},
	})

	mctest.PostTestResult()

}

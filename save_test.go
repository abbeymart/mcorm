// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mcorm

import (
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mctest"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
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
	//fmt.Printf("*****dbc-info: %v\n", dbc)
	// defer dbClose
	defer myDb.ClosePgxDbPool()
	// check db-connection-error
	if err != nil {
		fmt.Printf("*****db-connection-error: %v\n", err.Error())
		return
	}
	// expected db-connection result
	mcLogResult := mcauditlog.PgxLogParam{AuditDb: dbc.DbConn, AuditTable: TestAuditTable}
	// audit-log instance
	mcLog := mcauditlog.NewAuditLogPgx(dbc.DbConn, TestAuditTable)

	createCrudParams := types.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: CreateActionParams,
	}
	updateCrudParams := types.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: UpdateActionParams,
	}
	updateCrudParamsById := types.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: UpdateActionParamsById,
		RecordIds:    UpdateIds,
	}
	updateCrudParamsByParam := types.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: UpdateActionParamsByParam,
		QueryParams:  UpdateParams,
	}

	//fmt.Printf("test-action-params: %#v \n", createCrudParams.ActionParams)

	var crud interface{} = NewCrud(createCrudParams, TestCrudParamOptions)
	var updateCrud = NewCrud(updateCrudParams, TestCrudParamOptions)
	var updateIdCrud = NewCrud(updateCrudParamsById, TestCrudParamOptions)
	var updateParamCrud = NewCrud(updateCrudParamsByParam, TestCrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the Audit-DB and return an instance object:",
		TestFunc: func() {
			mctest.AssertEquals(t, err, nil, "error-response should be: nil")
			mctest.AssertEquals(t, mcLog, mcLogResult, "db-connection instance should be: "+mcLogResult.String())
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the CRUD-object and return an instance object:",
		TestFunc: func() {
			_, ok := crud.(*Crud)
			mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should create two new records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save([]string{})
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(types.CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should update two records and return success:",
		TestFunc: func() {
			res := updateCrud.Save(UpdateTableFields)
			fmt.Printf("updates: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by Ids and return success:",
		TestFunc: func() {
			res := updateIdCrud.Save(UpdateTableFields)
			fmt.Printf("update-by-ids: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by query-params and return success:",
		TestFunc: func() {
			res := updateParamCrud.Save([]string{})
			fmt.Printf("update-by-params: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params should return code: success")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should update two records, log-task and return success:",
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
			res := updateCrud.UpdateLog(updateCrud.ActionParams, GetTableFields, UpdateTableFields, tableFieldPointers)
			fmt.Printf("update-log: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-log should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by Ids, log-task and return success:",
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
			res := updateIdCrud.UpdateByIdLog(updateIdCrud.ActionParams, GetTableFields, UpdateTableFields, tableFieldPointers)
			fmt.Printf("update-by-ids-log: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id-log should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by query-params, log-task and return success:",
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
			res := updateParamCrud.UpdateByParamLog(updateParamCrud.ActionParams, GetTableFields, UpdateTableFields, tableFieldPointers)
			fmt.Printf("update-by-params-log: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params-log should return code: success")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should create two new records and return success[save-record-method]:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			var (
				id            string
				tableName     string
				logRecords    interface{}
				newLogRecords interface{}
				logBy         string
				logType       string
				logAt         time.Time
			)
			crud.RecordIds = []string{}
			crud.QueryParams = types.QueryParamType{}
			tableFieldPointers := []interface{}{&id, &tableName, &logRecords, &newLogRecords, &logBy, &logType, &logAt}
			// get-record method params
			saveRecParams := types.SaveCrudParamsType{
				CreateTableFields:  CreateTableFields,
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
			}
			res := crud.SaveRecord(saveRecParams)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(types.CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records and return success[save-record-method]:",
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
			updateCrud.RecordIds = []string{}
			updateCrud.QueryParams = types.QueryParamType{}
			// get-record method params
			saveRecParams := types.SaveCrudParamsType{
				UpdateTableFields:  UpdateTableFields,
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
				AuditLog:           true,
			}
			res := updateCrud.SaveRecord(saveRecParams)
			fmt.Printf("update[save-record]: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-log should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by Ids and return success[save-record-method]:",
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
			updateIdCrud.RecordIds = UpdateIds
			updateIdCrud.QueryParams = types.QueryParamType{}
			// get-record method params
			saveRecParams := types.SaveCrudParamsType{
				UpdateTableFields:  UpdateTableFields,
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
				AuditLog:           true,
			}
			res := updateIdCrud.SaveRecord(saveRecParams)
			fmt.Printf("update-by-ids[save-record]: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id-log should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by query-params and return success[save-record-method]:",
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
			updateParamCrud.RecordIds = []string{}
			updateParamCrud.QueryParams = UpdateParams
			// get-record method params
			saveRecParams := types.SaveCrudParamsType{
				UpdateTableFields:  UpdateTableFields,
				GetTableFields:     GetTableFields,
				TableFieldPointers: tableFieldPointers,
				AuditLog:           true,
			}
			res := updateParamCrud.SaveRecord(saveRecParams)
			fmt.Printf("update-by-params[save-record]: %#v \n", res)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params should return code: success")
		},
	})

	mctest.PostTestResult()

}

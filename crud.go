// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: Base type/function CRUD operations for PgDB

package mcorm

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/tasks"
	"github.com/abbeymart/mcresponse"
)

// Crud object / struct
type Crud struct {
	types.CrudParamsType
	types.CrudOptionsType
	CurrentRecords []interface{}
	TransLog       mcauditlog.PgxLogParam
	HashKey        string // Unique for exactly the same query
}

// NewCrud constructor returns a new crud-instance
func NewCrud(params types.CrudParamsType, options types.CrudOptionsType) (crudInstance *Crud) {
	crudInstance = &Crud{}
	// compute crud params
	crudInstance.AppDb = params.AppDb
	crudInstance.TableName = params.TableName
	crudInstance.UserInfo = params.UserInfo
	crudInstance.ActionParams = params.ActionParams
	crudInstance.RecordIds = params.RecordIds
	crudInstance.QueryParams = params.QueryParams
	crudInstance.SortParams = params.SortParams
	crudInstance.ProjectParams = params.ProjectParams
	crudInstance.ExistParams = params.ExistParams
	crudInstance.Token = params.Token
	crudInstance.TaskName = params.TaskName
	crudInstance.Skip = params.Skip
	crudInstance.Limit = params.Limit

	// crud options
	crudInstance.MaxQueryLimit = options.MaxQueryLimit
	crudInstance.AuditTable = options.AuditTable
	crudInstance.AccessTable = options.AccessTable
	crudInstance.RoleTable = options.RoleTable
	crudInstance.UserTable = options.UserTable
	crudInstance.UserProfileTable = options.UserProfileTable
	crudInstance.ServiceTable = options.ServiceTable
	crudInstance.AuditDb = options.AuditDb
	crudInstance.AccessDb = options.AccessDb
	crudInstance.LogCrud = options.LogCrud
	crudInstance.LogRead = options.LogRead
	crudInstance.LogCreate = options.LogCreate
	crudInstance.LogUpdate = options.LogUpdate
	crudInstance.LogDelete = options.LogDelete
	crudInstance.CheckAccess = options.CheckAccess // Dec 09/2020: user to implement auth as a middleware
	crudInstance.CacheExpire = options.CacheExpire // cache expire in secs
	// Compute HashKey from TableName, QueryParams, SortParams, ProjectParams and RecordIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.RecordIds)
	crudInstance.HashKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)

	// Default values
	if crudInstance.AuditTable == "" {
		crudInstance.AuditTable = "audits"
	}
	if crudInstance.AccessTable == "" {
		crudInstance.AccessTable = "access_keys"
	}
	if crudInstance.RoleTable == "" {
		crudInstance.RoleTable = "roles"
	}
	if crudInstance.UserTable == "" {
		crudInstance.UserTable = "users"
	}
	if crudInstance.UserProfileTable == "" {
		crudInstance.UserProfileTable = "user_profile"
	}
	if crudInstance.ServiceTable == "" {
		crudInstance.ServiceTable = "services"
	}
	if crudInstance.AuditDb == nil {
		crudInstance.AuditDb = crudInstance.AppDb
	}
	if crudInstance.AccessDb == nil {
		crudInstance.AccessDb = crudInstance.AppDb
	}
	if crudInstance.Skip < 0 {
		crudInstance.Skip = 0
	}

	if crudInstance.MaxQueryLimit == 0 {
		crudInstance.MaxQueryLimit = 10000
	}

	if crudInstance.Limit > crudInstance.MaxQueryLimit && crudInstance.MaxQueryLimit != 0 {
		crudInstance.Limit = crudInstance.MaxQueryLimit
	}

	if crudInstance.CacheExpire <= 0 {
		crudInstance.CacheExpire = 300 // 300 secs, 5 minutes
	}

	// Audit/TransLog instance
	crudInstance.TransLog = mcauditlog.NewAuditLogPgx(crudInstance.AuditDb, crudInstance.AuditTable)

	return crudInstance
}

// String() function implementation for crud instance/object
func (crud Crud) String() string {
	return fmt.Sprintf("CRUD Instance Information: %#v \n\n", crud)
}

// Methods

// SaveRecord function creates new record(s) or updates existing record(s)
func (crud *Crud) SaveRecord(params types.SaveCrudParamsType) mcresponse.ResponseMessage {
	//  compute taskType-records from actionParams: create or update
	var (
		createRecs types.ActionParamsType // records without id field-value
		updateRecs types.ActionParamsType // records with id field-value
		recIds     []string                 // capture recordIds for separate/multiple updates
	)
	for _, rec := range crud.ActionParams {
		// determine if record exists (update) or is new (create)
		if fieldValue, ok := rec["id"]; ok && fieldValue != "" {
			// validate fieldValue as string
			switch fieldValue.(type) {
			case string:
				updateRecs = append(updateRecs, rec)
				recIds = append(recIds, fieldValue.(string))
			default:
				// invalid fieldValue type (string)
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Invalid fieldValue type for fieldName: id, in record: %v", rec),
					Value:   nil,
				})
			}
		} else if len(crud.ActionParams) == 1 && (len(crud.RecordIds) > 0 || len(crud.QueryParams) > 0) {
			updateRecs = append(updateRecs, rec)
		} else {
			createRecs = append(createRecs, rec)
		}
	}

	// permit only create or update, not both at the same time
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "You may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}

	if len(createRecs) > 0 {
		// check task-permission - create
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(tasks.Create)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// save-record(s): create/insert new record(s): len(recordIds) = 0 && len(createRecs) > 0
		return crud.CreateBatch(createRecs, params.CreateTableFields)
	}

	// check task-permission - updates
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(tasks.Update)
		if accessRes.Code != "success" {
			return accessRes
		}
	}

	// update each record by it's recordId
	if len(updateRecs) >= 1 && (len(recIds) == len(updateRecs)) {
		if params.AuditLog || crud.LogUpdate {
			return crud.UpdateLog(updateRecs, params.GetTableFields, params.UpdateTableFields, params.TableFieldPointers)
		}
		return crud.Update(updateRecs, params.UpdateTableFields)
	}

	// update record(s) by recordIds
	if len(updateRecs) == 1 && len(crud.RecordIds) > 0 {
		if params.AuditLog || crud.LogUpdate {
			return crud.UpdateByIdLog(updateRecs, params.GetTableFields, params.UpdateTableFields, params.TableFieldPointers)
		}
		return crud.UpdateById(updateRecs, params.UpdateTableFields)
	}

	// update record(s) by queryParams
	if len(updateRecs) == 1 && len(crud.QueryParams) > 0 {
		if params.AuditLog || crud.LogUpdate {
			return crud.UpdateByParamLog(updateRecs, params.GetTableFields, params.UpdateTableFields, params.TableFieldPointers)
		}
		return crud.UpdateByParam(updateRecs, params.UpdateTableFields)
	}

	// otherwise return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})
}

// DeleteRecord function deletes/removes record(s) by id(s) or params
func (crud *Crud) DeleteRecord(params types.DeleteCrudParamsType) mcresponse.ResponseMessage {
	// check task-permission - delete
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(tasks.Delete)
		if accessRes.Code != "success" {
			return accessRes
		}
	}

	// delete-by-id
	if len(crud.RecordIds) > 0 {
		if params.AuditLog || crud.LogDelete {
			return crud.DeleteByIdLog(params.GetTableFields, params.TableFieldPointers)
		}
		return crud.DeleteById()
	}

	// delete-by-param
	if len(crud.QueryParams) > 0 {
		if params.AuditLog || crud.LogDelete {
			return crud.DeleteByParamLog(params.GetTableFields, params.TableFieldPointers)
		}
		return crud.DeleteByParam()
	}

	// delete-all ***RESTRICTED***

	// otherwise return error
	return mcresponse.GetResMessage("removeError", mcresponse.ResponseMessageOptions{
		Message: "Remove error: incomplete or invalid query-conditions provided",
		Value:   nil,
	})
}

// GetRecord function get records by id, params or all
func (crud *Crud) GetRecord(params types.GetCrudParamsType) mcresponse.ResponseMessage {
	// check task-permission - get/read
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(tasks.Read)
		if accessRes.Code != "success" {
			return accessRes
		}
	}

	// get-by-id
	if len(crud.RecordIds) > 0 {
		return crud.GetById(params.GetTableFields, params.TableFieldPointers)
	}

	// get-by-param
	if len(crud.QueryParams) > 0 {
		return crud.GetByParam(params.GetTableFields, params.TableFieldPointers)
	}

	// get-all-up-to-limit
	return crud.GetAll(params.GetTableFields, params.TableFieldPointers)
}

// GetRecords function get records by id, params or all - lookup-items
func (crud *Crud) GetRecords(params types.GetCrudParamsType) mcresponse.ResponseMessage {
	// get-by-id
	if len(crud.RecordIds) > 0 {
		return crud.GetById(params.GetTableFields, params.TableFieldPointers)
	}

	// get-by-param
	if len(crud.QueryParams) > 0 {
		return crud.GetByParam(params.GetTableFields, params.TableFieldPointers)
	}

	// get-all-up-to-limit
	return crud.GetAll(params.GetTableFields, params.TableFieldPointers)
}

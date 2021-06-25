// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: delete or remove record(s)

package mcorm

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccache"
	"github.com/abbeymart/mcorm/helper"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/tasks"
	"github.com/abbeymart/mcresponse"
)

// DeleteById method deletes or removes record(s) by record-id(s)
func (crud *Crud) DeleteById() mcresponse.ResponseMessage {
	// compute delete query by record-ids
	deleteQuery, dQErr := helper.ComputeDeleteQueryById(crud.TableName, crud.RecordIds)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully",
		Value:   commandTag.Delete(),
	})
}

// DeleteByParam method deletes or removes record(s) by query-parameters or where conditions
func (crud *Crud) DeleteByParam() mcresponse.ResponseMessage {
	// compute delete query by query-params
	deleteQuery, dQErr := helper.ComputeDeleteQueryByParam(crud.TableName, crud.QueryParams)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully",
		Value:   commandTag.Delete(),
	})
}

// DeleteAll method deletes or removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud *Crud) DeleteAll() mcresponse.ResponseMessage {
	// ***** perform DELETE-ALL-RECORDS FROM A TABLE, IF RELATIONS/CONSTRAINTS PERMIT *****
	// ***** && IF-AND-ONLY-IF-YOU-KNOW-WHAT-YOU-ARE-DOING *****
	// compute delete query
	delQuery := fmt.Sprintf("DELETE FROM %v", crud.TableName)
	commandTag, delErr := crud.AppDb.Exec(context.Background(), delQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// delete cache, by key (TableName)
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "key")

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: map[string]string{"query_desc": "all-records"},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully | " + logMessage,
		Value:   commandTag.Delete(),
	})
}

func (crud *Crud) DeleteByIdLog(tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// get records to delete, for audit-log
	if crud.LogDelete && len(tableFields) == len(tableFieldPointers) {
		getRes := crud.GetById(tableFields, tableFieldPointers)
		value, _ := getRes.Value.(types.CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}

	// perform delete-by-id
	delRes := crud.DeleteById()

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	// overall response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: delRes.Message + " | " + logMessage,
		Value:   delRes.Value,
	})
}

func (crud *Crud) DeleteByParamLog(tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// get records to delete, for audit-log
	if crud.LogDelete && len(tableFields) == len(tableFieldPointers) {
		getRes := crud.GetByParam(tableFields, tableFieldPointers)
		value, _ := getRes.Value.(types.CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}

	// perform delete-by-param
	delRes := crud.DeleteByParam()

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	// overall response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: delRes.Message + " | " + logMessage,
		Value:   delRes.Value,
	})
}

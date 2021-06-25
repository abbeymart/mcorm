// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"strings"
)

// ComputeSelectQueryAll compose select SQL script to retrieve all table-records
// The query may be limit/response may be controlled, by the user, by appending skip and limit options
func ComputeSelectQueryAll(tableName string, tableFields []string) (string, error) {
	if tableName == "" || len(tableFields) < 1 {
		return "", errors.New("table-name and table-fields are required to perform the select operation")
	}
	selectQuery := fmt.Sprintf("SELECT %v FROM %v", strings.Join(tableFields, ", "), tableName)
	return selectQuery, nil
}

// ComputeSelectQueryById compose select SQL script by id(s)
func ComputeSelectQueryById(tableName string, recordIds []string, tableFields []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 || len(tableFields) < 1 {
		return "", errors.New("table-name, table-fields and record-ids are required to perform the select operation")
	}
	// get record(s) based on projected/provided field names ([]string)
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	selectQuery += fmt.Sprintf("WHERE id IN( %v )", whereIds)
	return selectQuery, nil
}

// ComputeSelectQueryByParam compose SELECT query from the where-parameters
func ComputeSelectQueryByParam(tableName string, where types.QueryParamType, tableFields []string) (string, error) {
	if tableName == "" || len(where) < 1 || len(tableFields) < 1 {
		return "", errors.New("table-name, tableFields and where-params are required to perform the select operation")
	}
	// get record(s) based on projected/provided field names ([]string)
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	// add where-params condition
	if whereScript, err := ComputeWhereQuery(where); err == nil {
		selectQuery += whereScript
		return selectQuery, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}

// TODO: select-query functions for relational tables (eager & lazy queries) and data aggregation

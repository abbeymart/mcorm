// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute delete-SQL scripts

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
)

// ComputeDeleteQueryById function computes delete SQL script by id(s)
func ComputeDeleteQueryById(tableName string, recordIds []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 {
		return "", errors.New("table/collection name and doc-Ids are required for the delete-by-id operation")
	}
	// validated recordIds, strictly contains string/UUID values, to avoid SQL-injection
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	deleteQuery := "DELETE FROM " + tableName + " WHERE id IN(" + whereIds + ")"
	return deleteQuery, nil
}

// ComputeDeleteQueryByParam function computes delete SQL script by parameter specifications
func ComputeDeleteQueryByParam(tableName string, where types.QueryParamType) (string, error) {
	if tableName == "" || len(where) < 1 {
		return "", errors.New("table/collection name and where/query-condition are required for the delete-by-param operation")
	}
	if whereParam, err := ComputeWhereQuery(where); err == nil {
		deleteScript := fmt.Sprintf("DELETE FROM %v %v", tableName, whereParam)
		return deleteScript, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}

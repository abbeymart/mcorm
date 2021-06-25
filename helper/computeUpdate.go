// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute update-SQL scripts

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"github.com/asaskevich/govalidator"
	"time"
)

func ComputeUpdateQuery(tableName string, actionParams types.ActionParamsType, tableFields []string) ([]string, error) {
	if tableName == "" || len(actionParams) < 1 {
		return nil, errors.New("table-name and action-params are required for the update operation")
	}
	// compute tableFields from the first record, if len(tableFields) == 0
	if len(tableFields) == 0 {
		actRec := actionParams[0]
		for fName := range actRec {
			if fName == "id" {
				continue
			}
			tableFields = append(tableFields, fName)
		}
	}
	// compute update script from queryParams
	var updateQuery []string
	validUpdateItemCount := 0
	invalidUpdateItemCount := 0

	for recNum, rec := range actionParams {
		itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
		fieldCount := 0
		fieldLen := len(tableFields)
		for _, fieldName := range tableFields {
			fieldValue, ok := rec[fieldName]
			// check for the required fields in each record
			if !ok {
				return nil, errors.New(fmt.Sprintf("Record #%v [%#v]: required field_name[%v] is missing", recNum, rec, fieldName))
			}
			fieldCount += 1
			// update/set recFieldValues by fieldValue-type
			var currentFieldValue interface{}
			switch fieldValue.(type) {
			case time.Time:
				if fVal, ok := fieldValue.(time.Time); !ok {
					return nil, errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
				}
			case string:
				if fVal, ok := fieldValue.(string); !ok {
					return nil, errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					if govalidator.IsJSON(fVal) {
						if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
							return nil, errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
						} else {
							currentFieldValue = "'" + fValue + "'"
						}
					} else {
						currentFieldValue = "'" + fVal + "'"
					}
				}
			case int, uint, float64, bool:
				currentFieldValue = fieldValue
			default:
				// json-stringify fieldValue
				if fVal, err := json.Marshal(fieldValue); err != nil {
					return nil, errors.New(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
				} else {
					currentFieldValue = "'" + string(fVal) + "'"
				}
			}

			// add itemValue
			itemScript += fmt.Sprintf(" %v=%v", fieldName, currentFieldValue)
			if fieldLen > 1 && fieldCount < fieldLen {
				itemScript += ", "
			}
		}

		// add where condition by id
		itemScript += fmt.Sprintf(" WHERE id='%v'", rec["id"])
		//validate/update script content based on valid field specifications
		if fieldCount > 0 && fieldCount == fieldLen {
			validUpdateItemCount += 1
			updateQuery = append(updateQuery, itemScript)
		} else {
			invalidUpdateItemCount += 1
		}
	}
	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return nil, errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}
	return updateQuery, nil
}

func ComputeUpdateQueryById(tableName string, actionParams types.ActionParamsType, recordIds []string, tableFields []string) (string, error) {
	if tableName == "" || len(actionParams) < 1 || len(recordIds) < 1 {
		return "", errors.New("table-name, table-fields, action-params and record/doc-Ids are required for the update-by-id operation")
	}
	// compute tableFields from the first record, if len(tableFields) == 0
	if len(tableFields) == 0 {
		actRec := actionParams[0]
		for fName := range actRec {
			if fName == "id" {
				continue
			}
			tableFields = append(tableFields, fName)
		}
	}
	// compute update script from query-ids
	var updateQuery string
	itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	whereQuery := fmt.Sprintf(" WHERE id IN(%v)", whereIds)

	invalidUpdateItemCount := 0
	validUpdateItemCount := 0

	// only one actionParams record is required for update by docIds
	rec := actionParams[0]
	fieldCount := 0
	fieldLen := len(tableFields)
	for _, fieldName := range tableFields {
		fieldValue, ok := rec[fieldName]
		// check for the required fields in each record
		if !ok {
			return "", errors.New(fmt.Sprintf("Record [%#v]: required field_name[%v] is missing", rec, fieldName))
		}
		fieldCount += 1
		// update/set recFieldValues by fieldValue-type
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue = "'" + fValue + "'"
					}
				} else {
					currentFieldValue = "'" + fVal + "'"
				}
			}
		case int, uint, float64, bool:
			currentFieldValue = fieldValue
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return "", errors.New(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = "'" + string(fVal) + "'"
			}
		}
		// add itemValue
		itemScript += fmt.Sprintf(" %v=%v", fieldName, currentFieldValue)
		if fieldLen > 1 && fieldCount < fieldLen {
			itemScript += ", "
		}
	}
	//validate/update script content based on valid field specifications
	if fieldCount > 0 && fieldCount == fieldLen {
		validUpdateItemCount += 1
		updateQuery = itemScript + whereQuery
	} else {
		invalidUpdateItemCount += 1
	}

	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return "", errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}
	return updateQuery, nil
}

func ComputeUpdateQueryByParam(tableName string, actionParams types.ActionParamsType, where types.QueryParamType, tableFields []string) (string, error) {
	if tableName == "" || len(actionParams) < 1 || len(where) < 1 {
		return "", errors.New("table-name, action-params and where-params are required for the update-by-params operation")
	}
	// compute tableFields from the first record, if len(tableFields) == 0
	if len(tableFields) == 0 {
		actRec := actionParams[0]
		for fName := range actRec {
			if fName == "id" {
				continue
			}
			tableFields = append(tableFields, fName)
		}
	}

	// compute update script from queryParams
	var updateQuery string
	invalidUpdateItemCount := 0
	validUpdateItemCount := 0

	// only one actionParams record is required for update by where-params
	rec := actionParams[0]
	itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
	fieldCount := 0
	fieldLen := len(tableFields)
	for _, fieldName := range tableFields {
		fieldValue, ok := rec[fieldName]
		// check for the required fields in each record
		if !ok {
			return "", errors.New(fmt.Sprintf("Record [%#v]: required field_name[%v] is missing", rec, fieldName))
		}
		fieldCount += 1
		// update/set recFieldValues by fieldValue-type
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue = "'" + fValue + "'"
					}
				} else {
					currentFieldValue = "'" + fVal + "'"
				}
			}
		case int, uint, float64, bool:
			currentFieldValue = fieldValue
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return "", errors.New(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = "'" + string(fVal) + "'"
			}
		}
		// add itemValue
		itemScript += fmt.Sprintf(" %v=%v", fieldName, currentFieldValue)

		if fieldLen > 1 && fieldCount < fieldLen {
			itemScript += ", "
		}
	}
	//validate/update script content based on valid field specifications
	if fieldCount > 0 && fieldCount == fieldLen {
		validUpdateItemCount += 1
		updateQuery = itemScript
	} else {
		invalidUpdateItemCount += 1
	}

	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return "", errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}

	if whereScript, err := ComputeWhereQuery(where); err == nil {
		updateQuery += " " + whereScript
		return updateQuery, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}

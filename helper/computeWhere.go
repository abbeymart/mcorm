// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute where-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/groupOperators"
	"github.com/abbeymart/mcorm/types/operators"
	"sort"
	"strings"
	"time"
)

// ComputeWhereQuery function computes the multi-cases where-conditions for crud-operations
func ComputeWhereQuery(where types.QueryParamType) (string, error) {
	if len(where) < 1 {
		return "", errors.New("where condition is required")
	}
	// groups length/size
	groupsLen := len(where)
	// variable for empty group (i.e. contains no fieldItems)
	emptyGroupCount := 0
	// variable for valid group count, i.e. group with groupItems
	groupCount := 0
	// sort where by groupOrder (ASC)
	sort.SliceStable(where, func(i, j int) bool {
		return where[i].GroupOrder < where[j].GroupOrder
	})
	// compute where script from where
	// iterate through where (groups)
	whereQuery := "WHERE "
	for _, group := range where {
		var (
			emptyGroupItemCount = 0 // variable to determine empty/unspecified fieldName/Op/Value
			groupItemCount      = 0 // valid gItem count, i.e. group item with valid name and value
		)
		groupItemsLen := len(group.GroupItems) // total items in a group
		// check groupItems length, if 0 continue to the next group
		if groupItemsLen < 1 {
			emptyGroupCount += 1
			continue
		}
		// count valid group, i.e. group with groupItems
		groupCount += 1
		// sort group items by gItem/fieldOrder (ASC)
		gItems := group.GroupItems
		sort.SliceStable(gItems, func(i, j int) bool {
			return gItems[i].GroupItemOrder < gItems[j].GroupItemOrder
		})
		// compute the group-items query/script
		gItemQuery := "("
		for _, gItem := range gItems {
			// check gItem's fieldName, fieldOperator and fieldValue
			fieldName := ""
			fieldOperator := ""
			var fieldValue interface{}

			// ensure that len(gItem.GroupItem) == 1
			if len(gItem.GroupItem) != 1 {
				return "", errors.New(fmt.Sprintf("Only 1 field-name criteria is expected for each group-item"))
			}

			for fName, opVal := range gItem.GroupItem {
				fieldName = fName
				// ensure that len(opVal) == 1
				if len(opVal) != 1 {
					return "", errors.New(fmt.Sprintf("Only 1 operator-value criteria is expected for a field-name: %v", fieldName))
				}
				for fOp, val := range opVal {
					fieldOperator = fOp
					fieldValue = val
				}
			}
			if fieldName == "" || fieldOperator == "" || fieldValue == nil {
				// skip missing field/continue to the next gItem, or return error?
				emptyGroupItemCount += 1
				continue
				//return "", errors.New("field-name, operator and/or value are required")
			}
			// count valid gItem
			groupItemCount += 1
			switch strings.ToLower(fieldOperator) {
			case strings.ToLower(operators.Equals):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, currentFieldValue)
					}
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name:%v | field_value:%v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, "'"+fVal+"'")
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case []float64:
					if fVal, ok := fieldValue.([]float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.NotEquals):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, currentFieldValue)
					}
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, "'"+fVal+"'")
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<>%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.LessThan):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, currentFieldValue)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value: %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.LessThanOrEquals):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, currentFieldValue)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v<=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], type for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := gItem.GroupItemOp
				if groupItemOp == "" || !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = groupOperators.AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen-emptyGroupItemCount) {
					gItemQuery = gItemQuery + " " + strings.ToUpper(gItem.GroupItemOp) + " "
				}
			case strings.ToLower(operators.GreaterThan):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, currentFieldValue)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.GreaterThanOrEquals):
				switch fieldValue.(type) {
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						currentFieldValue := "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, currentFieldValue)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fmt.Sprintf("%v>=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.In):
				switch fieldValue.(type) {
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := ""
						idLen := len(fVal)
						for idCount, id := range fVal {
							inValues += "'" + id + "'"
							if idLen > 1 && idCount < idLen-1 {
								inValues += ", "
							}
						}
						gItemQuery += fmt.Sprintf("%v IN (%v)", fieldName, inValues)
					}
				case []bool:
					if fVal, ok := fieldValue.([]bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v IN %v", fieldName, inValues)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v IN %v", fieldName, inValues)
					}
				case []float32:
					if fVal, ok := fieldValue.([]float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v IN %v", fieldName, inValues)
					}
				case []float64:
					if fVal, ok := fieldValue.([]float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v IN %v", fieldName, inValues)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.NotIn):
				switch fieldValue.(type) {
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := ""
						idLen := len(fVal)
						for idCount, id := range fVal {
							inValues += "'" + id + "'"
							if idLen > 1 && idCount < idLen-1 {
								inValues += ", "
							}
						}
						gItemQuery += fmt.Sprintf("%v NOT IN (%v)", fieldName, inValues)
					}
				case []bool:
					if fVal, ok := fieldValue.([]bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v NOT IN %v", fieldName, inValues)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v NOT IN %v", fieldName, inValues)
					}
				case []float32:
					if fVal, ok := fieldValue.([]float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v NOT IN %v", fieldName, inValues)
					}
				case []float64:
					if fVal, ok := fieldValue.([]float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						inValues := "("
						fValLen := len(fVal)
						for i, v := range fVal {
							inValues += fmt.Sprintf("%v", v)
							if fValLen > 1 && i < fValLen-1 {
								inValues += inValues + ", "
							}
						}
						inValues += ")"
						gItemQuery += fmt.Sprintf("%v NOT IN %v", fieldName, inValues)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.StartsWith):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						//gItemQuery += fmt.Sprintf("%v LIKE %%v", fieldName, fVal)
						gItemQuery += fieldName + " LIKE " + fVal + "%"
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.EndsWith):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fieldName + " LIKE %" + fVal
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.NotStartsWith):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fieldName + " NOT LIKE " + fVal + "%"
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.NotEndsWith):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fieldName + " NOT LIKE %" + fVal
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.Includes):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fieldName + " LIKE %" + fVal + "%"
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			case strings.ToLower(operators.NotIncludes):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						gItemQuery += fieldName + " NOT LIKE %" + fVal + "%"
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field-name[%v] type for field-value %v", fieldName, fieldValue))
				}
			default:
				return "", errors.New(fmt.Sprintf("Unknown or unsupported field(%v) operator: %v", fieldName, fieldOperator))
			}
			//validate acceptable groupItemLinkOperators (and || or)
			gItemLinkOp := gItem.GroupItemOp
			gItemLinkOps := []string{"and", "or"}
			if gItemLinkOp == "" || !ArrayStringContains(gItemLinkOps, strings.ToLower(gItemLinkOp)) {
				gItemLinkOp = groupOperators.AND // default operator
			}
			// add gItemLinkOp, if groupItemsLen > 1 && groupItemCount < (groupItemsLen-emptyGroupItemCount)
			if groupItemsLen > 1 && groupItemCount < (groupItemsLen-emptyGroupItemCount) {
				gItemQuery += " " + strings.ToUpper(gItemLinkOp) + " "
			}
		}
		// continue to the next group iteration, if fieldItems is empty for the current group
		if emptyGroupItemCount == groupItemsLen {
			continue
		}
		// add closing bracket to complete the group-items query/script
		gItemQuery += ")"
		// add group-item-query to the group/where-script, append in sequence by groupOrder
		whereQuery += " " + gItemQuery
		//validate acceptable groupLinkOperators (and || or)
		grpLinkOp := group.GroupLinkOp
		groupLnOps := []string{"and", "or"}
		if grpLinkOp == "" || !ArrayStringContains(groupLnOps, strings.ToLower(grpLinkOp)) {
			grpLinkOp = groupOperators.AND // default operator
		}
		// add groupLinkOp, if groupsLen > 1 && groupCount < (groupsLen-emptyGroupCount)
		if groupsLen > 1 && groupCount < (groupsLen-emptyGroupCount) {
			whereQuery += " " + strings.ToUpper(grpLinkOp) + " "
		}
		// check WHERE script contains at least one condition, otherwise raise an exception
		if emptyGroupCount == groupsLen {
			return "", errors.New("no valid where condition specified")
		}
	}

	// if all went well, return valid where script
	return whereQuery, nil
}

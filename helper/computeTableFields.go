// @Author: abbeymart | Abi Akindele | @Created: 2020-12-15 | @Updated: 2020-12-15
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package helper

import (
	"errors"
	"github.com/abbeymart/mcorm/types"
)

func ComputeSaveFields(actionParams types.ActionParamsType, projectParams types.ProjectParamType) ([]string, error) {
	if len(actionParams) < 1 {
		return nil, errors.New("actionParams is required")
	}
	// obtain tableFields from api consumer (ProjectParams)
	var tableFields []string
	if len(projectParams) > 0 {
		for fieldName, ok := range projectParams {
			if ok {
				tableFields = append(tableFields, fieldName)
			}
		}
	}
	if len(tableFields) < 1 {
		// obtain tableFields from actionParams[0]
		for fieldName := range actionParams[0] {
			tableFields = append(tableFields, fieldName)
		}
	}

	return tableFields, nil
}

func ComputeGetFields(projectParams types.ProjectParamType) ([]string, error) {
	if len(projectParams) < 1 {
		return nil, errors.New("select/projection-params is required")
	}
	// obtain tableFields from api consumer (ProjectParams) | TODO: order-by-model-field-type-specs
	var tableFields []string
	if len(projectParams) > 0 {
		for fieldName, ok := range projectParams {
			if ok {
				tableFields = append(tableFields, fieldName)
			}
		}
		// include default fields (id) for get/select-query only
		// s = append([]int{0}, s...) | s = append([]string{"id"}, s...)
		if !ArrayStringContains(tableFields, "id") {
			tableFields = append([]string{"id"}, tableFields...)
		}
	}
	if len(tableFields) < 1 {
		return nil, errors.New("unable to compute query-fields")
	}

	return tableFields, nil
}


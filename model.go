// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: mcorm model

package mcorm

import (
	"fmt"
	"github.com/abbeymart/mcorm/helper"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/datatypes"
	"github.com/abbeymart/mcresponse"
	"github.com/asaskevich/govalidator"
	"strconv"
)

type CrudOperations interface {
	Save()
	Get()
	Delete()
}
type CrudSave interface {
	Save()
}
type CrudGet interface {
	Get()
}
type CrudDelete interface {
	Delete()
}

// Model object description
type Model struct {
	TaskName string
	types.ModelType
}

// NewModel constructor: for table structure definition
func NewModel(model types.ModelType) types.ModelType {
	result := types.ModelType{}
	result.AppDb = model.AppDb
	result.TableName = model.TableName
	result.RecordDesc = model.RecordDesc
	result.TimeStamp = model.TimeStamp
	result.ActorStamp = model.ActorStamp
	result.ActiveStamp = model.ActiveStamp
	result.Relations = model.Relations
	result.ComputedMethods = model.ComputedMethods
	result.ValidateMethods = model.ValidateMethods
	result.AlterSyncTable = model.AlterSyncTable

	// Default values
	if !result.TimeStamp {
		result.TimeStamp = true
	}
	if !result.ActiveStamp {
		result.ActiveStamp = true
	}
	if !result.ActorStamp {
		result.ActorStamp = true
	}

	return result
}

// GetParentRelations method computes the parent-relations for the current model table
func (model Model) GetParentRelations() []types.ModelRelationType {
	// extract relations/collections where targetTable == model-TableName
	var parentRelations []types.ModelRelationType
	modelRelations := model.Relations
	for _, item := range modelRelations {
		if item.TargetTable == model.TableName {
			parentRelations = append(parentRelations, item)
		}
	}
	return parentRelations
}

// GetChildRelations method computes the child-relations for the current model table
func (model Model) GetChildRelations() []types.ModelRelationType {
	// extract relations/collections where sourceTable == model-TableName
	var childRelations []types.ModelRelationType
	modelRelations := model.Relations
	for _, item := range modelRelations {
		if item.SourceTable == model.TableName {
			childRelations = append(childRelations, item)
		}
	}
	return childRelations
}

// GetParentTables method compose the parent-tables from GetParentRelations method response
func (model Model) GetParentTables() []string {
	var parentTables []string
	parentRelations := model.GetParentRelations()
	for _, rel := range parentRelations {
		parentTables = append(parentTables, rel.SourceTable)
	}
	return parentTables
}

// GetChildTables method compose the child-tables from GetParentRelations method response
func (model Model) GetChildTables() []string {
	var childTables []string
	childRelations := model.GetChildRelations()
	for _, rel := range childRelations {
		childTables = append(childTables, rel.TargetTable)
	}
	return childTables
}

// ComputeRecordValueType ComputeRecordValueType computes the corresponding standard/define types based on the record-fields types
func (model Model) ComputeRecordValueType(recordValue types.ActionParamType) types.ValueToDataType {
	computedType := types.ValueToDataType{}
	// perform computation of doc-value-types
	for key, val := range recordValue {
		// array check
		//if govalidator.IsType(val, "string") {}
		// switch fmt.Sprintf("%T", val)
		switch val.(type) {
		case []string:
			computedType[key] = datatypes.ArrayOfString
		case []int:
			computedType[key] = datatypes.ArrayOfNumber
		case []float32, []float64:
			computedType[key] = datatypes.ArrayOfNumber
		case []bool:
			computedType[key] = datatypes.ArrayOfBoolean
		//case []map:
		//	computedType[key] = datatypes.ArrayOfObject
		case []struct{}:
			computedType[key] = datatypes.ArrayOfStruct
		//case map:
		//	computedType[key] = datatypes.Map
		case struct{}:
			computedType[key] = datatypes.Object
		case string:
			// compute string value
			strVal := val.(string)
			//jsonStr, _ := json.Marshal(val)
			//strVal := string(jsonStr)
			var strToNum float64
			if val, err := strconv.Atoi(strVal); err == nil {
				strToNum = float64(val)
			}
			// check all string-based formats
			// TODO: ISO2, ISO3, Currency, Mime, JWT, PostalCode
			if govalidator.IsEmail(strVal) {
				computedType[key] = datatypes.Email
			} else if govalidator.IsUnixTime(strVal) {
				computedType[key] = datatypes.DateTime
			} else if govalidator.IsTime(strVal, "HH:MM:SS") {
				computedType[key] = datatypes.Time
			} else if govalidator.IsMongoID(strVal) {
				computedType[key] = datatypes.MongoDBId
			} else if govalidator.IsUUID(strVal) {
				computedType[key] = datatypes.UUID
			} else if govalidator.IsUUIDv3(strVal) {
				computedType[key] = datatypes.UUID3
			} else if govalidator.IsUUIDv4(strVal) {
				computedType[key] = datatypes.UUID4
			} else if govalidator.IsUUIDv5(strVal) {
				computedType[key] = datatypes.UUID5
			} else if govalidator.IsMD4(strVal) {
				computedType[key] = datatypes.MD4
			} else if govalidator.IsMD5(strVal) {
				computedType[key] = datatypes.MD5
			} else if govalidator.IsSHA1(strVal) {
				computedType[key] = datatypes.SHA1
			} else if govalidator.IsSHA256(strVal) {
				computedType[key] = datatypes.SHA256
			} else if govalidator.IsSHA384(strVal) {
				computedType[key] = datatypes.SHA384
			} else if govalidator.IsSHA512(strVal) {
				computedType[key] = datatypes.SHA512
			} else if govalidator.IsJSON(strVal) {
				computedType[key] = datatypes.JSON
			} else if govalidator.IsCreditCard(strVal) {
				computedType[key] = datatypes.CreditCard
			} else if govalidator.IsURL(strVal) {
				computedType[key] = datatypes.URL
			} else if govalidator.IsDNSName(strVal) {
				computedType[key] = datatypes.DomainName
			} else if govalidator.IsPort(strVal) {
				computedType[key] = datatypes.Port
			} else if govalidator.IsIP(strVal) {
				computedType[key] = datatypes.IP
			} else if govalidator.IsIPv4(strVal) {
				computedType[key] = datatypes.IP4
			} else if govalidator.IsIPv6(strVal) {
				computedType[key] = datatypes.IP6
			} else if govalidator.IsIMEI(strVal) {
				computedType[key] = datatypes.IMEI
			} else if govalidator.IsLatitude(strVal) {
				computedType[key] = datatypes.Latitude
			} else if govalidator.IsLongitude(strVal) {
				computedType[key] = datatypes.Longitude
			} else if govalidator.IsMAC(strVal) {
				computedType[key] = datatypes.MACAddress
			} else if govalidator.IsInt(strVal) {
				computedType[key] = datatypes.Integer
			} else if govalidator.IsPositive(strToNum) {
				computedType[key] = datatypes.Positive
			} else if govalidator.IsNegative(strToNum) {
				computedType[key] = datatypes.Negative
			} else if govalidator.IsNatural(strToNum) {
				computedType[key] = datatypes.Natural
			} else {
				computedType[key] = datatypes.String
			}
		case int, int8, int16, int32, int64:
			computedType[key] = datatypes.Integer
		case uint, uint8, uint16, uint32, uint64:
			computedType[key] = datatypes.Positive
		case float32, float64:
			computedType[key] = datatypes.Float
		case bool:
			computedType[key] = datatypes.Boolean
		default:
			computedType[key] = datatypes.Undefined
		}
	}
	return computedType
}

// UpdateDefaultValue method update default-value for non-null field with no specified value
// and pre-set value, prior to save (create/update) using setValueMethod
func (model Model) UpdateDefaultValue(recordValue types.ActionParamType) (setRecordValue types.ActionParamType) {
	// set default values, for null fields | then setValue (pre-set/transform), if specified
	// set base recordValue
	setRecordValue = recordValue
	// perform update of default/set-values for the doc-values => modelRecordValue
	for key, fieldValue := range recordValue {
		// defaultValue setting applies to FieldDescType only | otherwise, the value is required (not null)
		// transform fieldDesc to interface{} for type checking
		var fieldDescType interface{} = model.RecordDesc[key]
		//fieldValue := recordValue[key]
		// set default values
		if fieldValue != nil {
			switch fieldDescType.(type) {
			case types.FieldDescType:
				// type of defaultValue and fieldValue must be equivalent (re: validateMethod)
				fieldDesc := model.RecordDesc[key]
				if fieldDesc.DefaultValue != nil {
					defaultValue := fieldDesc.DefaultValue()
					// defaultValue and fieldValue types must match => validation-check
					// update setRecordValue for the key/field-column
					setRecordValue[key] = defaultValue
				}
			}
		}
		// setValue / transform field-value prior-to/before save-task (create / update)
		switch fieldDescType.(type) {
		case types.FieldDescType:
			setFieldValue := setRecordValue[key]
			if setFieldValue != nil && model.RecordDesc[key].SetValue != nil {
				// set/pre-set setRecordValue for the key/field-column
				setRecordValue[key] = model.RecordDesc[key].SetValue(recordValue)
			}
		}
	}
	return setRecordValue
}

// ValidateRecordValue method validate record-field-values based on model constraints and validation method
func (model Model) ValidateRecordValue(modelRecordValue types.ActionParamType, taskName string) types.ValidateResponseType {
	// perform validation of model-record-value
	// recommendation: use updated recordValue, defaultValues and setValues, prior to validation
	// get recordValue transformed types
	recordValueTypes := model.ComputeRecordValueType(modelRecordValue)
	// model-description/definition
	recordDesc := model.RecordDesc
	// combine errors/messages
	validateErrorMessage := map[string]string{}
	// perform model-recordValue validation
	for key, recordFieldValue := range modelRecordValue {
		// check field description / definition exists
		if recordFieldDesc, ok := recordDesc[key]; ok {
			// transform recordFieldDesc to interface{} for type checking
			var recordFieldDescType interface{} = recordFieldDesc
			switch recordFieldDescType.(type) {
			case types.FieldDescType:
				// validate fieldValue and fieldDesc (model) types
				// exception for fieldTypes: Text...
				typePermitted := recordValueTypes[key] == datatypes.String && recordFieldDesc.FieldType == datatypes.Text
				if recordValueTypes[key] != recordFieldDesc.FieldType && !typePermitted {
					errMsg := fmt.Sprintf("Invalid Type for:  %v. Expected %v, Got %v", key, recordFieldDesc.FieldType, recordValueTypes[key])
					if recordFieldDesc.ValidateMessage != "" {
						validateErrorMessage[key] = recordFieldDesc.ValidateMessage + " :: " + errMsg
					} else {
						validateErrorMessage[key] = errMsg
					}
				}

				// validate allowNull, fieldLength, min/maxValues...| user-defined-validation-methods
				// use values from transform docValue, including default/set-values
				// nullCheck, if recordField value is not specified
				if recordFieldValue == nil && !recordFieldDesc.AllowNull {
					errMsg := fmt.Sprintf("Value is required for: %v. Can't be Null", key)
					if recordFieldDesc.ValidateMessage != "" {
						validateErrorMessage[key+"-nullValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
					} else {
						validateErrorMessage[key+"-nullValidation"] = errMsg
					}
				}
				// validate field-value-type constraints: fieldLength, min/maxValues..
				switch recordFieldValue.(type) {
				case string:
					if fieldValue, ok := recordFieldValue.(string); ok {
						if recordFieldDesc.FieldLength > 0 {
							fieldLength := len(fieldValue)
							if fieldLength > recordFieldDesc.FieldLength {
								errMsg := fmt.Sprintf("Size of %v cannot be longer than %v", key, recordFieldDesc.FieldLength)
								if recordFieldDesc.ValidateMessage != "" {
									validateErrorMessage[key+"-lengthValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
								} else {
									validateErrorMessage[key+"-lengthValidation"] = errMsg
								}
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case int:
					if fieldValue, ok := recordFieldValue.(int); ok {
						if fieldValue < recordFieldDesc.MinValue && fieldValue > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v, and less than %v", key, recordFieldDesc.MinValue, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minMaxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minMaxValidation"] = errMsg
							}
						} else if fieldValue < recordFieldDesc.MinValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v", key, recordFieldDesc.MinValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minValidation"] = errMsg
							}
						} else if fieldValue > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be less than %v", key, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-maxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-maxValidation"] = errMsg
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case float32, float64:
					if fieldValue, ok := recordFieldValue.(float64); ok {
						if fieldValue < float64(recordFieldDesc.MinValue) && fieldValue > float64(recordFieldDesc.MaxValue) {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v, and less than %v", key, recordFieldDesc.MinValue, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minMaxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minMaxValidation"] = errMsg
							}
						} else if fieldValue < float64(recordFieldDesc.MinValue) {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v", key, recordFieldDesc.MinValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minValidation"] = errMsg
							}
						} else if fieldValue > float64(recordFieldDesc.MaxValue) {
							errMsg := fmt.Sprintf("Value of: %v must be less than %v", key, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-maxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-maxValidation"] = errMsg
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []string:
					if fieldValue, ok := recordFieldValue.([]string); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []int:
					if fieldValue, ok := recordFieldValue.([]int); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []float64, []float32:
					if fieldValue, ok := recordFieldValue.([]float64); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []struct{}:
					if fieldValue, ok := recordFieldValue.([]struct{}); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				}
			default:
				// validate field-value/type
				// use values from transform docValue, including default/set-values
				//if fieldValue, ok := modelRecordValue[key]; ok {
				//	fmt.Println(fieldValue)
				//}
			}
		} else {
			validateErrorMessage[key] = fmt.Sprintf("Invalid key: %v is not defined in the model", key)
		}
	}

	// perform user-defined recordValue validation
	// get validate method for the recordValue task by taskName (e.g. registerUser, login, saveProfile etc.)
	if modelValidateMethod, ok := model.ValidateMethods[taskName]; ok {
		valRes := modelValidateMethod(modelRecordValue)
		if !valRes.Ok {
			var modelErrorMsg = ""
			for _, msg := range valRes.Errors {
				if modelErrorMsg != "" {
					modelErrorMsg += " | " + msg
				} else {
					modelErrorMsg = msg
				}
			}
			validateErrorMessage[model.TableName+"-validationError"] = modelErrorMsg
		}
	}

	// check validateErrors
	if len(validateErrorMessage) != 0 {
		return types.ValidateResponseType{
			Ok:     false,
			Errors: validateErrorMessage,
		}
	}

	// return success, if validation process has been completed without errors
	var errMsg = types.MessageObject{}
	return types.ValidateResponseType{Ok: true, Errors: errMsg}
}

// Save method: sql.DB CRUD methods [pg, sqlite3...]
// Save method performs create (new records) or update (for current/existing records) task
func (model Model) Save(params types.CrudParamsType, options types.CrudOptionsType, tableFields []string) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName
	model.TaskName = params.TaskName
	if model.TaskName == "" {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: "taskName is required.",
			Value:   nil,
		})
	}
	// validate task/actionParams (recordValue), prior to saving, via model.ValidateRecordValue
	var actParams types.ActionParamsType
	if params.ActionParams != nil && len(params.ActionParams) > 0 {
		for _, recordValue := range params.ActionParams {
			// update defaultValues and setValues, before/prior to save
			modelRecordValue := model.UpdateDefaultValue(recordValue)
			// validate actionParam-item (recordValue) field-value
			validateRes := model.ValidateRecordValue(modelRecordValue, model.TaskName)
			if !validateRes.Ok || len(validateRes.Errors) > 0 {
				return helper.GetParamsMessage(validateRes.Errors)
			}
			// update actParams
			actParams = append(actParams, recordValue)
		}
	} else {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: "action-params is required to perform save operation.",
			Value:   nil,
		})
	}
	// TODO: update CRUD params and options
	params.ActionParams = actParams
	if !model.ActiveStamp {
		model.ActiveStamp = true
	}
	if !model.ActorStamp {
		model.ActorStamp = true
	}
	if !model.TimeStamp {
		model.TimeStamp = true
	}
	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform save-task
	return crud.Save(tableFields)
}

// Get method query the DB by record-id, defined query-parameter or all records, constrained
// by skip, limit and projected-field-parameters
func (model Model) Get(params types.CrudParamsType, options types.CrudOptionsType, tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName

	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform get-task
	return crud.GetById(tableFields, tableFieldPointers)
}

// GetStream method query the DB by record-ids, defined query-parameter or all records, constrained
// by skip, limit and projected-field-parameters, and stream the result
func (model Model) GetStream(params types.CrudParamsType, options types.CrudOptionsType) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName

	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform get-stream-task
	return crud.GetStream()
}

// DeleteById method delete record(s) by record-ids
func (model Model) DeleteById(params types.CrudParamsType, options types.CrudOptionsType) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName

	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform delete-task
	return crud.DeleteById()
}

// DeleteByParam method delete record(s) by specified query-parameter
func (model Model) DeleteByParam(params types.CrudParamsType, options types.CrudOptionsType) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName

	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform delete-task
	return crud.DeleteByParam()
}

// DeleteAll method delete all records from a table - ***** recommended for admin users only *****
func (model Model) DeleteAll(params types.CrudParamsType, options types.CrudOptionsType) mcresponse.ResponseMessage {
	// model specific params
	params.TableName = model.TableName

	// instantiate Crud action
	crud := NewCrud(params, options)
	// perform delete-task
	return crud.DeleteAll()
}

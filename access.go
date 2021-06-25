// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: optional access methods, to be used as middleware, prior to CRUD operation

package mcorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/abbeymart/mcorm/helper"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcorm/types/tasks"
	"github.com/abbeymart/mcresponse"
	"github.com/abbeymart/mctypes"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
)

// AccessInfoType for CheckUserAccess method value (interface{}) response,
// and to assert returned value
type AccessInfoType struct {
	UserId   string
	Group    string
	Groups   []string
	IsAdmin  bool
	IsActive bool
}

// TaskPermissionType for TaskPermission method value (interface{}) response,
// and to assert returned value
type TaskPermissionType struct {
	Ok       bool
	IsAdmin  bool
	IsActive bool
	UserId   string
	Group    string
	Groups   []string
}

// TaskPermission method determines the access permission by owner, role/group (on coll/table or doc/record(s)) or admin
// for various tasks: create/insert, update, delete/remove, read
func (crud *Crud) TaskPermission(taskType string) mcresponse.ResponseMessage {
	// permit crud tasks: by owner, role/group (on coll/table or doc/record(s)) or admin
	// task permission access variables
	var (
		taskPermitted   = false
		ownerPermitted  = false
		recordPermitted = false
		tablePermitted  = false
		isAdmin         = false
		isActive        = false
		userId          = ""
		tableId         = ""
		group           = ""
		groups          []string
		roleServices    []types.RoleServiceType
	)

	// check role-based access
	accessRes := crud.CheckTaskAccess()
	// capture roleServices value
	if accessRes.Code != "success" {
		return accessRes
	}

	// get access-record
	accessRec, ok := accessRes.Value.(types.CheckAccessType)
	if !ok {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Error parsing task access information/value",
			Value:   nil,
		})
	}
	// set access status variables
	isAdmin = accessRec.IsAdmin
	isActive = accessRec.IsActive
	roleServices = accessRec.RoleServices
	userId = accessRec.UserId
	group = accessRec.Group
	groups = accessRec.Groups
	tableId = accessRec.TableId

	// validate active status
	if !isActive {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Account is not active. Validate active status",
			Value:   nil,
		})
	}
	// validate task (roleServices) permission, for non-admin users
	if !isAdmin && len(roleServices) < 1 {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "You are not authorized to perform the requested action/task",
			Value:   nil,
		})
	}

	// determine records/documents ownership, for all records (atomic)
	accessUserId := accessRec.UserId
	recordIds := crud.RecordIds
	if len(recordIds) > 0 && accessUserId != "" && accessRec.IsActive {
		// SQL script
		sqlScript := fmt.Sprintf("SELECT id FROM %v WHERE id IN ($1) AND created_by = $2", crud.TableName)
		inValues := ""
		idLen := len(recordIds)
		for idCount, id := range recordIds {
			inValues += "'" + id + "'"
			if idLen > 1 && idCount < idLen-1 {
				inValues += ", "
			}
		}
		rows, err := crud.AppDb.Query(context.Background(), sqlScript, inValues, accessUserId)
		if err != nil {
			errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
		}
		defer rows.Close()
		// check rows count
		var rowCount = 0
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err == nil {
				rowCount += 1
			}
		}
		// ensure complete record count, as requested
		if rowCount == len(recordIds) {
			ownerPermitted = true
		}
	}

	// filter the roleServices by categories ("collection | table" or "record | document")
	collTabFunc := func(item types.RoleServiceType) bool {
		return item.ServiceCategory == tableId
	}
	recordFunc := func(item types.RoleServiceType) bool {
		return helper.ArrayStringContains(recordIds, item.ServiceCategory)
	}

	var (
		roleTables, roleRecords []types.RoleServiceType
	)
	if len(roleServices) > 0 {
		for _, v := range roleServices {
			if collTabFunc(v) {
				roleTables = append(roleTables, v)
			}
		}
		for _, v := range roleServices {
			if recordFunc(v) {
				roleRecords = append(roleRecords, v)
			}
		}
	}

	// helper functions
	canCreateFunc := func(item types.RoleServiceType) bool {
		return item.CanCreate
	}
	canUpdateFunc := func(item types.RoleServiceType) bool {
		return item.CanUpdate
	}
	canDeleteFunc := func(item types.RoleServiceType) bool {
		return item.CanDelete
	}
	canReadFunc := func(item types.RoleServiceType) bool {
		return item.CanRead
	}

	roleUpdateFunc := func(it1 string, it2 types.RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanUpdate
	}
	roleDeleteFunc := func(it1 string, it2 types.RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanDelete
	}
	roleReadFunc := func(it1 string, it2 types.RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanRead
	}

	roleRecFunc := func(it1 string, roleRecs []types.RoleServiceType, roleFunc types.RoleFuncType) bool {
		// test if any/some of the roleRecords it1/it2 met the access condition
		for _, it2 := range roleRecs {
			if roleFunc(it1, it2) {
				return true
			}
		}
		return false
	}

	// taskType specific permission(s)
	if !isAdmin && len(roleServices) > 0 {
		switch taskType {
		case tasks.Create, tasks.Insert:
			// collection/table level access | only tableId was included in serviceIds
			// must be able to perform create on the specified tableId(s)
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canCreateFunc(v) {
							return false
						}
					}
					return true
				}()
			}
		case tasks.Update:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canUpdateFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleUpdateFunc) {
							return false
						}
					}
					return true
				}()
			}
		case tasks.Delete, tasks.Remove:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canDeleteFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleDeleteFunc) {
							return false
						}
					}
					return true
				}()
			}
		case tasks.Read:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canReadFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleReadFunc) {
							return false
						}
					}
					return true
				}()
			}
		default:
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: "Unknown access type or access type not specified.",
				Value:   nil,
			})
		}
	}

	// overall access permitted
	taskPermitted = recordPermitted || tablePermitted || ownerPermitted || isAdmin

	if !taskPermitted {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "You are not authorized to perform the requested action/task.",
			Value:   TaskPermissionType{
				Ok: taskPermitted,
			},
		})
	}
	// const value = {...ok, ...{isAdmin, isActive, userId, group, groups}};
	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   TaskPermissionType{
			Ok:       taskPermitted,
			IsAdmin:  isAdmin,
			IsActive: isActive,
			UserId:   userId,
			Group:    group,
			Groups:   groups,
		},
	})
}

// CheckTaskAccess method determines the access by role-assignment
func (crud *Crud) CheckTaskAccess() mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status
	accessRes := crud.CheckUserAccess()
	if accessRes.Code != "success" {
		return accessRes
	}

	// set current-user info for next steps
	var (
		uId      string
		group    string
		groups   []string
		isAdmin  bool
		isActive bool
	)
	if val, ok := accessRes.Value.(AccessInfoType); ok {
		uId = val.UserId
		group = val.Group
		groups = val.Groups
		isAdmin = val.IsAdmin
		isActive = val.IsActive
	} else {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Error parsing user access information/value",
			Value:   nil,
		})
	}

	// if all the above checks passed, check for role-services access by taskType
	// obtain table/collName id(_id) from serviceTable/Coll (repo for all resources)
	var (
		serviceId string
		category  string
	)
	serviceScript := fmt.Sprintf("SELECT id, category from %v WHERE name=$1", crud.ServiceTable)
	serviceRow := crud.AccessDb.QueryRow(context.Background(), serviceScript, crud.TableName)
	// check error
	if err := serviceRow.Scan(&serviceId, &category); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Unauthorized: user information not found or inactive | %v", err.Error()),
			Value:   nil,
		})
	}
	// if permitted, include table/collId and recordIds in serviceIds
	tableId := ""
	serviceIds := crud.RecordIds
	catLowercase := strings.ToLower(category)
	if catLowercase == "table" || catLowercase == "collection" {
		tableId = serviceId
		serviceIds = append(serviceIds, serviceId)
	}

	var roleServices []types.RoleServiceType
	var rsErr error
	if len(serviceIds) > 0 {
		roleServices, rsErr = crud.GetRoleServices(crud.AccessDb, crud.RoleTable, group, serviceIds)
		if rsErr != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Action un-authorised / not-permitted | %v", rsErr.Error()),
				Value:   nil,
			})
		}
	}

	permittedRes := types.CheckAccessType{
		UserId:       uId,
		Group:        group,
		Groups:       groups,
		IsActive:     isActive,
		IsAdmin:      isAdmin,
		RoleServices: roleServices,
		TableId:      tableId,
	}

	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   permittedRes,
	})
}

// GetRoleServices method process and returns the permission to user / user-group for the specified service items
func (crud *Crud) GetRoleServices(accessDb *pgxpool.Pool, roleTable string, groupId string, serviceIds []string) ([]types.RoleServiceType, error) {
	var roleServices []types.RoleServiceType
	roleScript := fmt.Sprintf("SELECT id, service_id, service_category, can_read, can_create, can_delete, can_update from %v WHERE service_id IN ($1) AND group_id=$2 AND is_active=$3", roleTable)
	// where-in-values
	inValues := ""
	idLen := len(serviceIds)
	for idCount, id := range serviceIds {
		inValues += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			inValues += ", "
		}
	}

	rows, err := accessDb.Query(context.Background(), roleScript, inValues, groupId, true)
	if err != nil {
		//errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
		return roleServices, errors.New(fmt.Sprintf("%v", err.Error()))
	}
	defer rows.Close()
	var (
		roleId, serviceId, serviceCategory       string
		canRead, canCreate, canDelete, canUpdate bool
	)
	for rows.Next() {
		if err := rows.Scan(&roleId, &serviceId, &serviceCategory, &canRead, &canCreate, &canDelete, &canUpdate); err == nil {
			roleServices = append(roleServices, types.RoleServiceType{
				ServiceId:       serviceId,
				RoleId:          roleId,
				ServiceCategory: serviceCategory,
				CanRead:         canRead,
				CanCreate:       canCreate,
				CanUpdate:       canUpdate,
				CanDelete:       canDelete,
			})
		}
	}

	return roleServices, nil
}

// CheckUserAccess method determines the user access status: active, valid login and admin
func (crud *Crud) CheckUserAccess() mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status
	// get the accessKey information for the user
	accessScript := fmt.Sprintf("SELECT expire from %v WHERE user_id=$1 AND token=$2 AND login_name=$3", crud.AccessTable)
	rowAccess := crud.AccessDb.QueryRow(context.Background(), accessScript, crud.UserInfo.UserId, crud.UserInfo.Token, crud.UserInfo.LoginName)
	// check login-status/expiration
	var accessExpire int64
	if err := rowAccess.Scan(&accessExpire); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: please ensure that you are logged-in",
			Value:   nil,
		})
	} else {
		if (time.Now().Unix() * 1000) > accessExpire {
			return mcresponse.GetResMessage("tokenExpired", mcresponse.ResponseMessageOptions{
				Message: "Access expired: please login to continue",
				Value:   nil,
			})
		}
	}
	// check the current-user status/info
	var (
		uId      string
		group    string
		groups   []string
		isAdmin  bool
		isActive bool
	)
	userScript := fmt.Sprintf("SELECT id, groups, isAdmin, isActive from %v WHERE id=$1 AND is_active=$2", crud.UserTable)
	rowUser := crud.AccessDb.QueryRow(context.Background(), userScript, crud.UserInfo.UserId, true)
	if err := rowUser.Scan(&uId, &groups, &isAdmin, &isActive); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: user information not found or is inactive",
			Value:   nil,
		})
	}
	// get default-group from user profile
	pScript := fmt.Sprintf("SELECT group from %v WHERE user_id=$1 is_active=$2", crud.UserProfileTable)
	userProfile := crud.AccessDb.QueryRow(context.Background(), pScript, crud.UserInfo.UserId, true)
	if err := userProfile.Scan(&group); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: user-profile-group information not found or is inactive",
			Value:   nil,
		})
	}

	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value: AccessInfoType{
			UserId:   uId,
			Group:    group,
			Groups:   groups,
			IsAdmin:  isAdmin,
			IsActive: isActive,
		},
	})
}

// CheckLoginStatus method checks if the user exists and has active login status/token
func (crud *Crud) CheckLoginStatus(params mctypes.UserInfoType) mcresponse.ResponseMessage {
	// check if user exists, from users table
	emailUsername := helper.EmailUsername(params.LoginName)
	email := emailUsername.Email
	username := emailUsername.Username
	var uId string
	if email != "" {
		query := fmt.Sprintf("SELECT id from $1 WHERE id=$2 AND email=$3")
		row := crud.AccessDb.QueryRow(context.Background(), query, crud.UserTable, params.UserId, email)
		err := row.Scan(&uId)
		if err != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Record not found for %v. Register a new account", params.LoginName),
				Value:   nil,
			})
		}
	} else if username != "" {
		query := fmt.Sprintf("SELECT id from $1 WHERE id=$2 AND username=$3")
		row := crud.AccessDb.QueryRow(context.Background(), query, crud.UserTable, params.UserId, username)
		err := row.Scan(&uId)
		if err != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Record not found for %v. Register a new account", params.LoginName),
				Value:   nil,
			})
		}
	} else {
		// invalid user-information provided
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Invalid user-information provided.",
			Value:   nil,
		})
	}

	// check loginName, userId and token validity... from access_keys table
	var expire int64
	query := fmt.Sprintf("SELECT expire from $1 WHERE id=$2 AND login_name=$3 AND token=$4")
	row := crud.AccessDb.QueryRow(context.Background(), query, crud.AccessTable, params.UserId, params.LoginName, params.Token)
	err := row.Scan(&expire)
	if err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Access information for %v not found. Login first, or contact system administrator", params.LoginName),
			Value:   nil,
		})
	}
	if (time.Now().Unix() * 1000) > expire {
		// Delete the expired access_keys | remove access-info from access_keys table
		delQuery := fmt.Sprintf("DELETE FROM %v WHERE id=$1 AND token=$2", crud.AccessTable)
		_, _ = crud.AppDb.Exec(context.Background(), delQuery, params.UserId, params.Token)
		return mcresponse.GetResMessage("tokenExpired", mcresponse.ResponseMessageOptions{
			Message: "Access expired: please login to continue",
			Value:   nil,
		})
	}

	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / Access permitted.",
		Value:   uId,
	})
}

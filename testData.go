// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mcorm

import (
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcorm/helper"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mctypes"
	"time"
)

const TestTable = "audits_test1"
const DeleteAllTable = "audits_test2"
const TestAuditTable = "audits"

var CreateTableFields = []string{
	"table_name",
	"log_records",
	"log_type",
	"log_by",
	"log_at",
}

var UpdateTableFields = []string{
	"table_name",
	"log_records",
	"log_type",
	"log_at",
}

var DeleteSelectTableFields = []string{
	"id",
	"table_name",
	"log_records",
	"new_log_records",
	"log_by",
	"log_type",
	"log_at",
}

var GetTableFields = []string{
	"id",
	"table_name",
	"log_records",
	"new_log_records",
	"log_by",
	"log_type",
	"log_at",
}

type TestParam struct {
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Url      string  `json:"url"`
	Priority int     `json:"priority"`
	Cost     float64 `json:"cost"`
}

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de"

var TestUserInfo = mctypes.UserInfoType{
	UserId:    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	LoginName: "abbeymart",
	Email:     "abbeya1@yahoo.com",
	Language:  "en-US",
	FirstName: "Abi",
	LastName:  "Akindele",
	Token:     "",
	Expire:    0,
	Group:     "TBD",
}

var Recs = TestParam{Name: "Abi", Desc: "Testing only", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var TableRecords, _ = helper.DataToValueParam(Recs)

// NewRecs fmt.Println("table-records-json", string(tableRecords))
var NewRecs = TestParam{Name: "Abi Akindele", Desc: "Testing only - updated", Url: "localhost:9900", Priority: 1, Cost: 2000.00}
var NewTableRecords, _ = helper.DataToValueParam(NewRecs)

//fmt.Println("new-table-records-json", string(newTableRecords))
//var ReadP = map[string][]string{"keywords": {"lagos", "nigeria", "ghana", "accra"}}
//var ReadParams, _ = json.Marshal(ReadP)

var TestCrudParamOptions = types.CrudOptionsType{
	AuditTable:    "audits",
	UserTable:     "users",
	ServiceTable:  "services",
	AccessTable:   "access_keys",
	VerifyTable:   "verify_users",
	RoleTable:     "roles",
	LogCreate:     true,
	LogUpdate:     true,
	LogDelete:     true,
	LogRead:       true,
	LogLogin:      true,
	LogLogout:     true,
	MaxQueryLimit: 100000,
	MsgFrom:       "support@mconnect.biz",
}

// CreateRecordA create record(s)
var CreateRecordA = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: TableRecords,
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}
var CreateRecordB = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: TableRecords,
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}
var valParam1, _ = helper.DataToValueParam(CreateRecordA)
var valParam2, _ = helper.DataToValueParam(CreateRecordB)
var CreateActionParams = types.ActionParamsType{
	valParam1,
	valParam2,
}

// UpdateRecordType update record(s)
type UpdateRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

var upRecs = TestParam{Name: "Abi100", Desc: "Testing only100", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords, _ = helper.DataToValueParam(upRecs)
var upRecs2 = TestParam{Name: "Abi200", Desc: "Testing only200", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords2, _ = helper.DataToValueParam(upRecs2)
var UpdateRecordA = UpdateRecordType{
	Id:            "d46a29db-a9a3-47b9-9598-e17a7338e474",
	TableName:     "services",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}
var UpdateRecordB = UpdateRecordType{
	Id:            "8fcdc5d5-f4e3-4f98-ba19-16e798f81070",
	TableName:     "services2",
	LogRecords:    upTableRecords2,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordById = UpdateRecordType{
	TableName:     "services2",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordByParam = mcauditlog.AuditRecord{
	TableName:     "services3",
	LogRecords:    upTableRecords2,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var UpdateParams = types.QueryParamType{
	types.QueryGroupType{
		GroupName:   "id_logtype",
		GroupOrder:  1,
		GroupLinkOp: "or",
		GroupItems: []types.QueryItemType{
			{
				GroupItem:      map[string]map[string]interface{}{"id": {"eq": "57d58438-2941-40f2-8e6f-c9e4539dab3e"}},
				GroupItemOrder: 1,
				GroupItemOp:    "and",
			},
			{
				GroupItem:      map[string]map[string]interface{}{"log_type": {"eq": "create"}},
				GroupItemOrder: 2,
				GroupItemOp:    "and",
			},
		},
	},
}

var updateRec1, _ = helper.DataToValueParam(UpdateRecordA)
var updateRec2, _ = helper.DataToValueParam(UpdateRecordB)
var updateRecId, _ = helper.DataToValueParam(UpdateRecordById)
var updateRecParam, _ = helper.DataToValueParam(UpdateRecordByParam)

var UpdateActionParams = types.ActionParamsType{
	updateRec1,
	updateRec2,
}

var UpdateActionParamsById = types.ActionParamsType{
	updateRecId,
}
var UpdateActionParamsByParam = types.ActionParamsType{
	updateRecParam,
}

// GetRecordType get record(s)
type GetRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

// GetIds get by ids & params
var GetIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GetParams = types.QueryParamType{
	types.QueryGroupType{
		GroupName:   "id_table",
		GroupOrder:  2,
		GroupLinkOp: "and",
		GroupItems: []types.QueryItemType{
			{
				GroupItem:      map[string]map[string]interface{}{"id": {"in": []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}}},
				GroupItemOrder: 1,
				GroupItemOp:    "and",
			},
			{
				GroupItem:      map[string]map[string]interface{}{"table_name": {"eq": "services"}},
				GroupItemOrder: 2,
				GroupItemOp:    "and",
			},
		},
	},
}

// DeleteIds delete record(s) by ids & params
var DeleteIds = []string{"dba4adbb-4482-4f3d-bb05-0db80c30876b", "02f83bc1-8fa3-432a-8432-709f0df3f3b0"}
var DeleteParams = types.QueryParamType{
	types.QueryGroupType{
		GroupName:   "id_table",
		GroupOrder:  2,
		GroupLinkOp: "and",
		GroupItems: []types.QueryItemType{
			{
				GroupItem:      map[string]map[string]interface{}{"id": {"eq": "57d58438-2941-40f2-8e6f-c9e4539dab3e"}},
				GroupItemOrder: 1,
				GroupItemOp:    "and",
			},
			{
				GroupItem:      map[string]map[string]interface{}{"table_name": {"eq": "services"}},
				GroupItemOrder: 2,
				GroupItemOp:    "and",
			},
		},
	},
	types.QueryGroupType{
		GroupName:   "id_logtype",
		GroupOrder:  1,
		GroupLinkOp: "or",
		GroupItems: []types.QueryItemType{
			{
				GroupItem:      map[string]map[string]interface{}{"id": {"eq": "57d58438-2941-40f2-8e6f-c9e4539dab3e"}},
				GroupItemOrder: 1,
				GroupItemOp:    "and",
			},
			{
				GroupItem:      map[string]map[string]interface{}{"log_type": {"eq": "create"}},
				GroupItemOrder: 2,
				GroupItemOp:    "and",
			},
		},
	},
}

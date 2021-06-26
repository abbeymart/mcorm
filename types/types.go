// @Author: abbeymart | Abi Akindele | @Created: 2020-12-22 | @Updated: 2020-12-22
// @Company: mConnect.biz | @License: MIT
// @Description: crud operations' types - updated

package types

import (
	"fmt"
	"github.com/abbeymart/mctypes"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleServiceType struct {
	ServiceId            string `json:"serviceId"`
	RoleId               string `json:"roleId"`
	ServiceCategory      string `json:"serviceCategory"`
	CanRead              bool   `json:"canRead"`
	CanCreate            bool   `json:"canCreate"`
	CanUpdate            bool   `json:"canUpdate"`
	CanDelete            bool   `json:"canDelete"`
	CanCrud              bool   `json:"canCrud"`
	TableAccessPermitted bool   `json:"tableAccessPermitted"`
}

type CheckAccessType struct {
	UserId       string            `json:"userId" mcorm:"userId"`
	Group        string            `json:"group" mcorm:"group"`
	Groups       []string          `json:"groups" mcorm:"groups"`
	IsActive     bool              `json:"isActive" mcorm:"isActive"`
	IsAdmin      bool              `json:"isAdmin" mcorm:"isAdmin"`
	RoleServices []RoleServiceType `json:"roleServices" mcorm:"roleServices"`
	TableId      string            `json:"tableId" mcorm:"tableId"`
}

type CheckAccessParamsType struct {
	AccessDb         *pgxpool.Pool        `json:"accessDb"`
	UserInfo         mctypes.UserInfoType `json:"userInfo"`
	TableName        string               `json:"tableName"`
	RecordIds        []string             `json:"recordIds"` // for update, delete and read tasks
	AccessTable      string               `json:"accessTable"`
	UserTable        string               `json:"userTable"`
	RoleTable        string               `json:"roleTable"`
	ServiceTable     string               `json:"serviceTable"`
	UserProfileTable string               `json:"userProfileTable"`
}

type RoleFuncType func(it1 string, it2 RoleServiceType) bool
type FieldValueType interface{}
type ActionParamType map[string]interface{}
type ValueToDataType map[string]interface{}
type ActionParamsType []ActionParamType
type ExistParamType map[string]interface{}
type ExistParamsType []ExistParamType
type SortParamType map[string]int     // 1 for "asc", -1 for "desc"
type ProjectParamType map[string]bool // 1 or true for inclusion, 0 or false for exclusion

type QueryItemType struct {
	GroupItem      map[string]map[string]interface{} `json:"groupItem"`      // key1 => fieldName, key2 => fieldOperator, interface{}=> value(s)
	GroupItemOrder int                               `json:"groupItemOrder"` // item/field order within the group
	GroupItemOp    string                            `json:"groupItemOp"`    // group-item relationship to the next item (AND, OR), the last item groupItemOp should be "" or will be ignored
}

type QueryGroupType struct {
	GroupName   string          `json:"groupName"`   // for group-items(fields) categorization
	GroupItems  []QueryItemType `json:"groupItems"`  // group items to be composed by category
	GroupOrder  int             `json:"groupOrder"`  // group order
	GroupLinkOp string          `json:"groupLinkOp"` // group relationship to the next group (AND, OR), the last group groupLinkOp should be "" or will be ignored
}

type QueryParamType []QueryGroupType
type WhereParamType []QueryGroupType

// CrudParamsType is the struct type for receiving, composing and passing CRUD inputs
type CrudParamsType struct {
	AppDb         *pgxpool.Pool        `json:"-"`
	TableName     string               `json:"-"`
	UserInfo      mctypes.UserInfoType `json:"userInfo"`
	ActionParams  ActionParamsType     `json:"actionParams"`
	ExistParams   ExistParamsType      `json:"existParams"`
	QueryParams   QueryParamType       `json:"queryParams"`
	RecordIds     []string             `json:"recordIds"`
	ProjectParams ProjectParamType     `json:"projectParams"`
	SortParams    SortParamType        `json:"sortParams"`
	Token         string               `json:"token"`
	Skip          int                  `json:"skip"`
	Limit         int                  `json:"limit"`
	TaskType      string               `json:"-"`
}

type CrudOptionsType struct {
	ParentTables          []string
	ChildTables           []string
	RecursiveDelete       bool
	CheckAccess           bool
	AccessDb              *pgxpool.Pool
	AuditDb               *pgxpool.Pool
	ServiceDb             *pgxpool.Pool
	AuditTable            string
	ServiceTable          string
	UserTable             string
	RoleTable             string
	AccessTable           string
	VerifyTable           string
	UserProfileTable      string
	MaxQueryLimit         int
	LogCrud               bool
	LogCreate             bool
	LogUpdate             bool
	LogRead               bool
	LogDelete             bool
	LogLogin              bool
	LogLogout             bool
	UnAuthorizedMessage   string
	RecExistMessage       string
	CacheExpire           int
	LoginTimeout          int
	UsernameExistsMessage string
	EmailExistsMessage    string
	MsgFrom               string
}

type CrudParamType struct {
	AppDb            *pgxpool.Pool // use *pgxpool.Pool, preferred || *pgx.Conn
	TableName        string
	Token            string
	UserInfo         mctypes.UserInfoType
	UserId           string
	Group            string
	Groups           []string
	RecordIds        []string
	ActionParams     ActionParamsType
	QueryParams      QueryParamType
	ExistParams      ExistParamsType
	ProjectParams    ProjectParamType
	SortParams       SortParamType
	Skip             int
	Limit            int
	ParentTables     []string
	ChildTables      []string
	RecursiveDelete  bool
	CheckAccess      bool
	AccessDb         *pgxpool.Pool
	AuditDb          *pgxpool.Pool
	AuditTable       string
	ServiceTable     string
	UserTable        string
	UserProfileTable string
	RoleTable        string
	AccessTable      string
	MaxQueryLimit    int
	logCrud          bool
	LogCreate        bool
	LogUpdate        bool
	LogRead          bool
	LogDelete        bool
	//transLog AuditLog
	HashKey             string
	IsRecExist          bool
	ActionAuthorized    bool
	UnAuthorizedMessage string
	RecExistMessage     string
	IsAdmin             bool
	CreateItems         ActionParamsType
	UpdateItems         ActionParamsType
	CurrentRecs         ActionParamsType
	RoleServices        []RoleServiceType
	SubItems            []bool
	CacheExpire         int
	Params              CrudParamsType
}

// MongoDB specific types

type MongoCrudTaskType struct {
	AppDb         *mongo.Client
	TableName     string
	UserInfo      mctypes.UserInfoType
	ActionParams  ActionParamsType
	ExistParams   ExistParamsType
	QueryParams   QueryParamType
	RecordIds     []string
	ProjectParams ProjectParamType
	SortParams    SortParamType
	Token         string
	Options       MongoCrudOptionsType
	TaskName      string
}

type MongoCrudOptionsType struct {
	Skip                  int
	Limit                 int
	ParentTables          []string
	ChildTables           []string
	RecursiveDelete       bool
	CheckAccess           bool
	AccessDb              *mongo.Client
	AuditDb               *mongo.Client
	ServiceDb             *mongo.Client
	AuditTable            string
	ServiceTable          string
	UserTable             string
	RoleTable             string
	AccessTable           string
	VerifyTable           string
	MaxQueryLimit         int
	LogCrud               bool
	LogCreate             bool
	LogUpdate             bool
	LogRead               bool
	LogDelete             bool
	LogLogin              bool
	LogLogout             bool
	UnAuthorizedMessage   string
	RecExistMessage       string
	CacheExpire           int
	LoginTimeout          int
	UsernameExistsMessage string
	EmailExistsMessage    string
	MsgFrom               string
}

type MongoCrudParamType struct {
	AppDb           *mongo.Client
	TableName       string
	Token           string
	UserInfo        mctypes.UserInfoType
	UserId          string
	Group           string
	Groups          []string
	RecordIds       []string
	ActionParams    ActionParamsType
	QueryParams     QueryParamType
	ExistParams     ExistParamsType
	ProjectParams   ProjectParamType
	SortParams      SortParamType
	Skip            int
	Limit           int
	ParentTables    []string
	ChildTables     []string
	RecursiveDelete bool
	CheckAccess     bool
	AccessDb        *mongo.Client
	AuditDb         *mongo.Client
	AuditTable      string
	ServiceTable    string
	UserTable       string
	RoleTable       string
	AccessTable     string
	MaxQueryLimit   int
	logCrud         bool
	LogCreate       bool
	LogUpdate       bool
	LogRead         bool
	LogDelete       bool
	//transLog AuditLog
	HashKey             string
	IsRecExist          bool
	ActionAuthorized    bool
	UnAuthorizedMessage string
	RecExistMessage     string
	IsAdmin             bool
	CreateItems         ActionParamsType
	UpdateItems         ActionParamsType
	CurrentRecs         ActionParamsType
	RoleServices        []RoleServiceType
	SubItems            []bool
	CacheExpire         int
	Params              MongoCrudTaskType
}

type MessageObject map[string]string

type ValidateResponseType struct {
	Ok     bool          `json:"ok"`
	Errors MessageObject `json:"errors"`
}
type OkResponse struct {
	Ok bool `json:"ok"`
}

// CRUD operations

type CreateQueryResponseType struct {
	CreateQuery string
	FieldNames  []string
	FieldValues [][]interface{}
}

type UpdateQueryResponseType struct {
	UpdateQuery string
	WhereQuery  string
	FieldValues []interface{}
}

type WhereQueryResponseType struct {
	WhereQuery  string
	FieldValues []interface{}
}

type DeleteQueryResponseType struct {
	DeleteQuery string
	WhereQuery  string
	FieldValues []interface{}
}

type SelectQueryResponseType struct {
	SelectQuery string
	WhereQuery  string
	FieldValues []interface{}
}

type SaveParamsType struct {
	UserInfo    mctypes.UserInfoType `json:"userInfo"`
	QueryParams QueryParamType       `json:"queryParams"`
	RecordIds   []string             `json:"recordIds"`
	//ActionParams ActionParamsType `json:"actionParams"`
}

type DeleteParamsType struct {
	UserInfo    mctypes.UserInfoType `json:"userInfo"`
	RecordIds   []string             `json:"recordIds"`
	QueryParams QueryParamType       `json:"queryParams"`
}

type GetParamsType struct {
	UserInfo     mctypes.UserInfoType `json:"userInfo"`
	Skip         int                  `json:"skip"`
	Limit        int                  `json:"limit"`
	RecordIds    []string             `json:"recordIds"`
	QueryParams  QueryParamType       `json:"queryParams"`
	SortParam    SortParamType        `json:"sortParams"`
	ProjectParam ProjectParamType     `json:"projectParam"`
}

type SaveCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	CreateTableFields  []string
	UpdateTableFields  []string
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}

type DeleteCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}

type GetCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}

// ErrorType provides the structure for error reporting
type ErrorType struct {
	Code    string
	Message string
}

type SaveError ErrorType
type CreateError ErrorType
type UpdateError ErrorType
type DeleteError ErrorType
type ReadError ErrorType
type AuthError ErrorType
type ConnectError ErrorType
type SelectQueryError ErrorType
type WhereQueryError ErrorType
type CreateQueryError ErrorType
type UpdateQueryError ErrorType
type DeleteQueryError ErrorType

// sample Error() implementation
func (err ErrorType) Error() string {
	return fmt.Sprintf("Error-code: %v | Error-message: %v", err.Code, err.Message)
}

type CrudResultType struct {
	QueryParam   QueryParamType `json:"queryParam"`
	RecordIds    []string       `json:"recordIds"`
	RecordCount  int            `json:"recordCount"`
	TableRecords []interface{}  `json:"tableRecords"`
}

type LogRecordsType struct {
	TableFields  []string       `json:"tableFields"`
	TableRecords []interface{}  `json:"tableRecords"`
	QueryParam   QueryParamType `json:"queryParam"`
	RecordIds    []string       `json:"recordIds"`
}

// ORM types

type RecordValueType map[string]ActionParamType
type RecordDescType map[string]FieldDescType

type GetValueType func() interface{}
type SetValueType func(val interface{}) interface{}
type DefaultValueType func() interface{}
type ValidateMethodType func(val interface{}) bool
type ValidateMethodResponseType func(val interface{}) ValidateResponseType
type ComputedValueType func(val interface{}) interface{}

type ValidateMethodsType map[string]ValidateMethodResponseType
type ComputedMethodsType map[string]ComputedValueType

type FieldDescType struct {
	FieldType       string
	FieldLength     int    // default: 255 for DataType.STRING
	FieldPattern    string // "/^[0-9]{10}$/" => includes 10 digits, 0 to 9 | "/^[0-9]{6}.[0-9]{2}$/ => max 16 digits and 2 decimal places
	AllowNull       bool   // default: true
	Unique          bool
	Indexable       bool
	PrimaryKey      bool
	MinValue        int
	MaxValue        int
	SetValue        SetValueType       // set/transform fieldValue prior to save(create/insert), T=>fieldType
	DefaultValue    DefaultValueType   // result/T must be of fieldType
	Validate        ValidateMethodType // T=>fieldType, returns a bool (valid=true/invalid=false)
	ValidateMessage string
	Comments        string
}

type ModelRelationType struct {
	SourceTable   string
	TargetTable   string
	SourceField   string
	TargetField   string
	RelationType  string
	SourceModel   ModelType
	TargetModel   ModelType
	ForeignField  string // source-to-targetField map
	RelationField string // relation-targetField, for many-to-many
	RelationTable string // optional tableName for many-to-many | default: sourceTable_targetTable
	OnDelete      string
	OnUpdate      string
}

type ModelType struct {
	AppDb           *pgxpool.Pool
	TableName       string
	RecordDesc      RecordDescType
	TimeStamp       bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp      bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp     bool // record active status, isActive (true | false) | default: true
	Relations       []ModelRelationType
	ComputedMethods ComputedMethodsType // model-level functions, e.g fullName(a, b: T): T
	ValidateMethods ValidateMethodsType
	AlterSyncTable  bool // create / alter table/collection and sync existing data, if there was a change to the table structure | default: true
	// if alterSyncTable: false it will create/re-create the table, with no data sync
}

type UniqueFieldsType [][]string

type ModelOptionsType struct {
	TimeStamp       bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp      bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp     bool // auto-add isActive, if not already set | default: true
	RecordValueDesc RecordDescType
	RecordValue     ActionParamType
}

type ModelCrudOptionsType struct {
	ModelOptionsType
	Relations      []ModelRelationType // for ref-integrity
	UniqueFields   UniqueFieldsType    // composite unique-fields
	PrimaryFields  []string            // composite primary-fields
	RequiredFields []string            // may be computed from FieldDesc allowNull attributes
}

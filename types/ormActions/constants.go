// @Author: abbeymart | Abi Akindele | @Created: 2020-12-23 | @Updated: 2020-12-23
// @Company: mConnect.biz | @License: MIT
// @Description: ORM Table Relation Actions

package ormActions

const (
	Restrict = "restrict" // must remove target-record(s), prior to removing source-record
	Cascade  = "cascade"  // default for ON UPDATE | update foreignKey value or delete foreignKey record/value
	NoAction = "noaction" // leave the foreignKey value, as-is
	Default  = "default"  // set foreignKey to specified default value
	Null     = "null"     // set foreignKey value to null
)

// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute alter-table script

package helper

import (
	"github.com/abbeymart/mcorm/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateAlterTableQuery(model types.ModelType) (string, error)  {


	return "", nil
}

func CreateAlterTable(model types.ModelType, appDb *pgxpool.Pool) error {


	return nil
}

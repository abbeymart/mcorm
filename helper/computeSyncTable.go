// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute sync-table script / data migration actions

package helper

import (
	"github.com/abbeymart/mcorm/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SyncTableQuery(model types.ModelType, appDb *pgxpool.Pool) (string, error)  {


	return "", nil
}

func SyncTable(model types.ModelType, appDb *pgxpool.Pool) error {


	return nil
}
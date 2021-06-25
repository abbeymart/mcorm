// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute temporary-table script

package helper

import (
	"github.com/abbeymart/mcorm/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateTempTableQuery(model types.ModelType, appDb *pgxpool.Pool) (string, error) {


	return "", nil
}

func CreateTempTable(model types.ModelType, appDb *pgxpool.Pool) error {


	return nil
}

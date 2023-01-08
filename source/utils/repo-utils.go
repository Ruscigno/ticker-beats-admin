package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// OverTransactionAdditionalQueries allow you to runs additional taks over the openned transaction
type OverTransactionAdditionalQueries func(con *pgxpool.Conn) error

// ExecScript runs que SQL
func ExecScript(ctx context.Context, model interface{}, sql string, dbCon *sqlx.DB, ignoreFields []string) (int64, error) {
	var query string
	var values []interface{}
	var fields []string
	if model != nil {
		values, fields = StructToSlice(model, ignoreFields)
		params := []string{}
		for i := 0; i < len(fields); i++ {
			params = append(params, "$"+strconv.Itoa(i+1))
		}
		query = fmt.Sprintf(sql, strings.Join(fields[:], ","), strings.Join(params[:], ","))
	} else {
		query = sql
	}
	tx, err := dbCon.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	rs, err := dbCon.Exec(query, values...)
	if err != nil {
		zap.L().Error(err.Error(),
			zap.String("query", query),
			zap.String("values", spew.Sdump(values)),
		)
		return 0, err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return ra, nil
}

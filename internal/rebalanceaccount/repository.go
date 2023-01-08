package rebalanceaccount

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
)

type rebalanceAccountRepo struct {
	ctx   *context.Context
	dbCon *sqlx.DB
}

// NewRepository creates a new Repository
func NewRepository(ctx *context.Context, dbCon *sqlx.DB) RebalanceAccountRepo {
	return &rebalanceAccountRepo{ctx: ctx, dbCon: dbCon}
}

func (x *rebalanceAccountRepo) ChangeActiveState(sourceAccountID, destinationAccountID int64, active bool) error {
	const updateSignal = `
		UPDATE tickerbeats.signals
		SET active=$1
		WHERE sourceaccountid=$2 AND destinationaccountid=$3;
	`
	_, err := x.dbCon.Exec(updateSignal, active, sourceAccountID, destinationAccountID)
	if err != nil {
		return err
	}
	return nil
}

// RebalanceAccount rebalance the account
func (x *rebalanceAccountRepo) TopProfitableAccounts(limit int) ([]*AccountsToRebalanceModel, error) {
	const SelectQuery string = `
		SELECT ap.accountid, sum(profit) as profit, avg(profit_rate) as profit_rate, sum(volume) as volume, 
			sum(deals) as deals, sum(ap.recordid) as week_num, sum(profit*recordid)/6 as pond_profit
		FROM tickerbeats.accounts_profit ap,
			tickerbeats.accounts a
		where a.accountid = ap.accountid
		group by ap.accountid
		having sum(ap.recordid) = 6
		order by pond_profit desc
		limit $1
	`
	result := []*AccountsToRebalanceModel{}
	err := x.dbCon.Select(&result, SelectQuery, limit)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

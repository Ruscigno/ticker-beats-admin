package rebalanceaccount

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ruscigno/ticker-beats-admin/source/utils"
	"github.com/Ruscigno/ticker-heart/source/tickerbeats/signal"
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

func (x *rebalanceAccountRepo) ChangeActiveStateByID(signalID []string, active bool) error {
	var updateSignal = `
		UPDATE tickerbeats.signals
		SET active=$1
		WHERE signalID in (%s);
	`
	updateSignal = fmt.Sprintf(updateSignal, strings.Join(signalID, ","))
	_, err := x.dbCon.Exec(updateSignal, active)
	if err != nil {
		return err
	}
	return nil
}

func (x *rebalanceAccountRepo) ChangeActiveStateByDestinationAccountID(accountID int64, active bool) error {
	const updateSignal = `
	UPDATE tickerbeats.signals
	SET active=$1
	WHERE destinationaccountid=$2;
`
	_, err := x.dbCon.Exec(updateSignal, active, accountID)
	if err != nil {
		return err
	}
	return nil
}

// RebalanceAccount rebalance the account
func (x *rebalanceAccountRepo) TopProfitableAccounts(limit int) ([]*AccountsToRebalanceModel, error) {
	const SelectQuery string = `
		SELECT ap.accountid as account_id, sum(profit) as profit, avg(profit_rate) as profit_rate, sum(volume) as volume, 
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

func (x *rebalanceAccountRepo) TempInsertSignal(sourceAccountID, destinationAccountID int64) error {
	const insertSignal = `
	INSERT INTO tickerbeats.signals
		(signalid, sourceaccountid, destinationaccountid, active, maxdepositpercent, stopiflessthan, maxspread, minutestoexpire, orderboost, orderboosttype, deviation, created)
		VALUES(nextval('tickerbeats.signal_id'), $1, $2, true, 100, 0, 0, 2, 0.03, 4, 10, $3);
	`
	_, err := x.dbCon.Exec(insertSignal, sourceAccountID, destinationAccountID, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (x *rebalanceAccountRepo) TempGetAllSignals(onlyActive bool) ([]signal.Signal, error) {
	var SelectQuery string = "select %s from tickerbeats.signals where signalid > 0%s"
	active := ""
	if onlyActive {
		active = " and active = true"
	}
	signals := []signal.Signal{}
	_, fields := utils.StructToSlice(signal.Signal{}, nil)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","), active)
	err := x.dbCon.Select(&signals, query)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil
		}
		return nil, err
	}
	return signals, nil
}

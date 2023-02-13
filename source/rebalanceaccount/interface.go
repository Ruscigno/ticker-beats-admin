package rebalanceaccount

import (
	"github.com/Ruscigno/ticker-heart/source/tickerbeats/signal"
	"github.com/gin-gonic/gin"
)

type RebalanceAccount interface {
	RebalanceAccount(c *gin.Context, req RebalanceAccountRequest) error
}

type RebalanceAccountService interface {
	RebalanceAccount
}

type RebalanceAccountRepo interface {
	TopProfitableAccounts(limit int) ([]*AccountsToRebalanceModel, error)
	ChangeActiveState(sourceAccountID, destinationAccountID int64, active bool) error
	ChangeActiveStateByID(signalID []string, active bool) error
	ChangeActiveStateByDestinationAccountID(accountID int64, active bool) error
	TempInsertSignal(sourceAccountID, destinationAccountID int64) error
	TempGetAllSignals(aonlyActive bool) ([]signal.Signal, error)
}

package rebalanceaccount

import "github.com/gin-gonic/gin"

type RebalanceAccount interface {
	RebalanceAccount(c *gin.Context, req RebalanceAccountRequest) error
}

type RebalanceAccountService interface {
	RebalanceAccount
}

type RebalanceAccountRepo interface {
	TopProfitableAccounts(limit int) ([]*AccountsToRebalanceModel, error)
	ChangeActiveState(sourceAccountID, destinationAccountID int64, active bool) error
}

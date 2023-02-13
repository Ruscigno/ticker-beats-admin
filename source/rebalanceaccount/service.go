package rebalanceaccount

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rebalanceAccountService struct {
	ctx  *context.Context
	repo RebalanceAccountRepo
}

// NewRebalanceAccountService creates a new Service
func NewRebalanceAccountService(ctx *context.Context, repo RebalanceAccountRepo) RebalanceAccountService {
	return &rebalanceAccountService{
		ctx:  ctx,
		repo: repo,
	}
}

// RebalanceAccount rebalance the account
func (s *rebalanceAccountService) RebalanceAccount(c *gin.Context, req RebalanceAccountRequest) error {
	profitables, err := s.repo.TopProfitableAccounts(req.Limit)
	if err != nil {
		zap.L().Error("error getting top profitable accounts", zap.Error(err))
		return err
	}
	profitablesList := []int64{}
	for _, profitable := range profitables {
		profitablesList = append(profitablesList, profitable.AccountID)
	}
	for _, destinationAccountID := range req.AccountsID {
		err = s.repo.ChangeActiveStateByDestinationAccountID(destinationAccountID, false)
		if err != nil {
			zap.L().Error("error de-activating signals", zap.Any("destinationAccountID", destinationAccountID), zap.Error(err))
			return err
		}
		for _, sourceAccountID := range profitablesList {
			err = s.repo.TempInsertSignal(sourceAccountID, destinationAccountID)
			if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				zap.L().Error("error inserting signal", zap.Error(err))
				return err
			}
			err = s.repo.ChangeActiveState(sourceAccountID, destinationAccountID, true)
			if err != nil {
				zap.L().Error("error activating signal",
					zap.Int64("sourceAccountID", sourceAccountID),
					zap.Int64("destinationAccountID", destinationAccountID),
					zap.Error(err))
				return err
			}
		}

	}
	return nil
}

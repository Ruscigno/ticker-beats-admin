package rebalanceaccount

import (
	"context"

	ss "github.com/Ruscigno/ticker-heart/source/tickerbeats/signal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rebalanceAccountService struct {
	ctx    *context.Context
	repo   RebalanceAccountRepo
	signal ss.SignalService
}

// NewRebalanceAccountService creates a new Service
func NewRebalanceAccountService(ctx *context.Context, repo RebalanceAccountRepo, signal ss.SignalService) RebalanceAccountService {
	return &rebalanceAccountService{
		ctx:    ctx,
		repo:   repo,
		signal: signal,
	}
}

// RebalanceAccount rebalance the account
func (s *rebalanceAccountService) RebalanceAccount(c *gin.Context, req RebalanceAccountRequest) error {
	profitables, err := s.repo.TopProfitableAccounts(req.Limit)
	if err != nil {
		zap.L().Error("error getting top profitable accounts", zap.Error(err))
		return err
	}
	profitablesMap := make(map[int64]*AccountsToRebalanceModel)
	for _, profitable := range profitables {
		profitablesMap[profitable.AccountID] = profitable
	}
	signals, err := s.signal.GetSignalByDestinationAccountID(req.AccoundID, false)
	if err != nil {
		zap.L().Error("error getting signals by destination", zap.Error(err))
		return err
	}
	for _, signal := range signals {
		_, ok := profitablesMap[signal.SourceAccountID]
		err = s.repo.ChangeActiveState(signal.SourceAccountID, signal.DestinationAccountID, ok)
		if err != nil {
			zap.L().Error("error updating signal", zap.Error(err))
			return err
		}
	}
	return nil
}

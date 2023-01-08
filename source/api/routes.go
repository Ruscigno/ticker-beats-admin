package api

import (
	"context"

	"github.com/Ruscigno/ticker-beats-admin/source/rebalanceaccount"
	"github.com/gin-gonic/gin"
)

type ginRoutes struct {
	ctx *context.Context
	rb  rebalanceaccount.RebalanceAccountService
}

var ginRoutesInstance *ginRoutes

func GetRoutes(ctx *context.Context, rb rebalanceaccount.RebalanceAccountService) *gin.Engine {
	ginRoutesInstance = &ginRoutes{
		ctx: ctx,
		rb:  rb,
	}
	router := gin.Default()
	router.POST("/experts/rebalance-account", ginRoutesInstance.RebalanceAccount)
	return router
}

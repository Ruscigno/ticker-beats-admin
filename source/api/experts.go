package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Ruscigno/ticker-beats-admin/source/rebalanceaccount"
)

// e.POST(/experts/rebalance-account", RebalanceAccount)
func (g *ginRoutes) RebalanceAccount(c *gin.Context) {
	var req rebalanceaccount.RebalanceAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("error binding request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := g.rb.RebalanceAccount(c, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Return a JSON response
	c.Status(http.StatusOK)
}

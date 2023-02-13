package main

import (
	"context"
	"os"
	"time"

	"github.com/Ruscigno/ticker-beats-admin/source/api"
	"github.com/Ruscigno/ticker-beats-admin/source/rebalanceaccount"
	"github.com/Ruscigno/ticker-beats-admin/source/utils"
	"github.com/Ruscigno/ticker-beats-admin/source/utils/app"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
)

func main() {
	logger := app.SetupLogger()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	loc, err := time.LoadLocation("")
	if err != nil {
		zap.L().Fatal("unable to set timezone properly", zap.Error(err))
	}
	time.Local = loc // -> this is setting the global timezone

	cfg := &utils.Config{
		Database: utils.Database{
			DSN:                os.Getenv(utils.ConfigDatabaseDSN),
			MaxOpenConnections: os.Getenv(utils.ConfigDbMaxOpenConnections),
			MaxIdleConnections: os.Getenv(utils.ConfigDbIdleConnections),
		},
	}
	db, err := app.InitDatabase(cfg.Database)
	if err != nil {
		zap.L().Fatal("error connecting to database", zap.Error(err))
	}
	defer db.Close()

	ctx := context.Background()
	repo := rebalanceaccount.NewRepository(&ctx, db)
	serv := rebalanceaccount.NewRebalanceAccountService(&ctx, repo)

	r := api.GetRoutes(&ctx, serv)
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	err = r.Run(":31034")
	logger.Fatal("Internal server error", zap.Error(err))
}

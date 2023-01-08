package main

import (
	"context"
	"time"

	"os"

	"github.com/Ruscigno/ticker-beats-admin/internal/api"
	"github.com/Ruscigno/ticker-beats-admin/internal/utils"
	"github.com/Ruscigno/ticker-beats-admin/internal/utils/app"

	_ "github.com/lib/pq"

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
	router := api.GetRoutes(&ctx, nil)
	err = router.Run("localhost:31034")
	logger.Fatal("Internal server error", zap.Error(err))
}

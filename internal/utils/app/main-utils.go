package app

import (
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Ruscigno/ticker-beats-admin/internal/utils"

	_ "github.com/lib/pq"
)

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}

func GetWriteSyncer(logName string) zapcore.WriteSyncer {
	var ioWriter = &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    20, // MB
		MaxBackups: 5,  // number of backups
		MaxAge:     28, //days
		LocalTime:  true,
		Compress:   false, // disabled by default
	}
	var sw = WriteSyncer{
		ioWriter,
	}
	return sw
}

func SetupLogger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	var config zap.Config
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config = zap.NewDevelopmentConfig()
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	configConsole := config
	configConsole.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(configConsole.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func InitDatabase(dbCfg utils.Database) (*sqlx.DB, error) {
	var err error
	max, idle := 10, 10
	if dbCfg.MaxOpenConnections != "" {
		max, err = strconv.Atoi(dbCfg.MaxOpenConnections)
		if err != nil {
			return nil, err
		}
	}

	if dbCfg.MaxIdleConnections != "" {
		idle, err = strconv.Atoi(dbCfg.MaxIdleConnections)
		if err != nil {
			return nil, err
		}
	}
	db, err := sqlx.Connect("postgres", dbCfg.DSN)
	if err != nil {
		zap.L().Fatal("Error trying to connect to the database", zap.Error(err))
	}
	db.SetMaxOpenConns(max)
	db.SetMaxIdleConns(idle)
	db.SetConnMaxLifetime(time.Minute)
	return db, err
}

type Controllers struct {
	DB *sqlx.DB
	// TrSvc   *trService.TradeRulesService
	// InfoSvc *infoS.AccountsInfoService
	// AccSvc  *accService.AccountsService
	// SigSvc  *sigS.SignalService
	// Beats   *bb.TickerBeatsService
	// PosSvc  *posService.PositionsService
	// PoHSvc  *poHService.PositionsHistoryService
	// DeaSvc  *deaService.DealsService
	// OrdSvc  *ordService.OrdersService
	// TtSvc   *ttService.TradeTransactionService
	// SySvc   sy.ISymbolTranslationService
}

func InitControllers(ctx context.Context, db *sqlx.DB) *Controllers {
	// Repos
	// trRepo := trRepository.NewTradeRulesRepo(ctx, db)
	// infoRepo := infoR.NewAccountInfoRepo(ctx, db)
	// accRepo := accRepository.NewAccountsRepo(ctx, db)
	// sigRepo := sigR.NewSignalRepository(ctx, db)
	// posRepo := posRepository.NewPositionsRepo(ctx, db)
	// poHRepo := poHRepository.NewPositionsHistoryRepo(ctx, db)
	// deaRepo := deaRepository.NewDealsRepo(ctx, db)
	// ordRepo := ordRepository.NewOrdersRepo(ctx, db)
	// ttRepo := ttRepository.NewTradeTransactionRepo(ctx, db)
	// syRepo := sy.NewSymbolTranslationRepo(ctx, db)

	//Controllers
	// sySvc := sy.NewSymbolTranslationService(ctx, syRepo)
	// trSvc := trService.NewTradeRulesService(ctx, trRepo)
	// infoSvc := infoS.NewAccountsInfoService(ctx, infoRepo)
	// accSvc := accService.NewAccountsService(ctx, accRepo, infoRepo, infoSvc)
	// sigSvc := sigS.NewSignalService(ctx, sigRepo)
	// beats := bb.NewTickerBeatsService(ctx, sigSvc, trSvc, infoSvc, sySvc)
	// poHSvc := poHService.NewPositionsHistoryService(ctx, poHRepo)
	// posSvc := posService.NewPositionsService(ctx, poHSvc, posRepo, beats)
	// deaSvc := deaService.NewDealsService(ctx, deaRepo, posSvc, beats)
	// ordSvc := ordService.NewOrdersService(ctx, ordRepo, beats, deaSvc)
	// ttSvc := ttService.NewTradeTransactionService(ctx, ttRepo, beats)

	return &Controllers{
		DB: db,
		// TrSvc:   &trSvc,
		// InfoSvc: &infoSvc,
		// AccSvc:  &accSvc,
		// SigSvc:  &sigSvc,
		// Beats:   &beats,
		// PosSvc:  &posSvc,
		// PoHSvc:  &poHSvc,
		// DeaSvc:  &deaSvc,
		// OrdSvc:  &ordSvc,
		// TtSvc:   &ttSvc,
		// SySvc:   sySvc,
	}
}

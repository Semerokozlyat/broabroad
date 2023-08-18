package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/exp/slog"

	"broabroad/internal/app/api"
	"broabroad/internal/app/config"
	"broabroad/internal/app/database"
	"broabroad/internal/app/logic"
)

type MainApp struct {
	apiServer *api.Server
}

func newMainApp(cfg *config.Config, logger *slog.Logger) *MainApp {
	ctx := context.Background()
	logger.Info("initializing storage", slog.String("conn_string", cfg.DB.ConnectionString))
	dbRepository, err := database.NewRepo(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("init db repository: %v", err)
	}
	storageController := logic.NewStorageController(dbRepository)

	logger.Info("initializing API server", slog.String("address", cfg.HTTPServer.Address))
	srv := api.NewServer(cfg, storageController)

	return &MainApp{
		apiServer: srv,
	}
}

func (app *MainApp) Run(ctx context.Context) error {
	go func() {
		err := app.apiServer.Run()
		if err != nil {
			log.Fatalf("run API server: %s", err)
		}
	}()
	log.Printf("main app started")
	<-ctx.Done()
	log.Printf("main app stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := app.apiServer.Shutdown(ctxShutdown)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("graceful server shutdown failed: %v", err)
	}
	log.Printf("main app finished graceful exit")

	return nil
}

func setupLogger(logCfg config.Log) (*slog.Logger, error) {
	var level slog.Level
	switch logCfg.Level {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	default:
		return nil, fmt.Errorf("unsupported log level: %s", logCfg.Level)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	return logger, nil
}

func RunMainApp(ctx context.Context, cfg *config.Config) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		log.Printf("caught system call: %v", <-c)
		cancel()
	}()

	logger, err := setupLogger(cfg.Log)
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}
	app := newMainApp(cfg, logger)
	return app.Run(ctx)
}

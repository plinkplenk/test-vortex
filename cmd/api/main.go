package main

import (
	"errors"
	"flag"
	apiApp "github.com/plinkplenk/test-vortex/internal/app"
	"github.com/plinkplenk/test-vortex/internal/config"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var debug = flag.Bool("debug", false, "-debug")

func main() {
	flag.Parse()
	cfg := config.Setup()
	loggingLevel := slog.LevelInfo
	if *debug {
		loggingLevel = slog.LevelDebug
	}
	logger := setupLogger(loggingLevel)
	app, err := apiApp.New(apiApp.Params{Config: cfg, Logger: logger, Debug: *debug})
	if err != nil {
		panic(err)
	}

	go func() {
		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Error while running app", "error:", err)
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	logger.Info("Shutting down server...")
	app.Stop()
}

func setupLogger(level slog.Level) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: level,
			},
		),
	)
}

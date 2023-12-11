package main

import (
	"log/slog"
	"os"

	"github.com/jbockle/captivated/server"
	"github.com/jbockle/captivated/server/config"
	"github.com/jbockle/captivated/server/services"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("Exiting due to panic", "err", err)
			os.Exit(1)
		}
	}()

	setDefaultLogger()

	config.Init()
	services.Init()
	go services.StartDeleteExpiredTask()

	server.Serve()
}

func setDefaultLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // TODO move to env var
	}

	var handler slog.Handler
	// TODO switch handler based on env var
	handler = slog.NewTextHandler(os.Stdout, opts)
	// handler = slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
